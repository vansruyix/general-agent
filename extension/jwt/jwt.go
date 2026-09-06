package jwt

import (
	"errors"
	"fmt"
	"general-agent/extension/errorx"
	"general-agent/extension/logz"
	"time"

	"github.com/golang-jwt/jwt"
)

var TokenExpireDuration = time.Hour * 1
var TokenRefreshDuration = time.Minute * 10

type MyClaims struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	MachineID string `json:"machine_id" binding:"required"`
	AssetID   string `json:"asset_id" binding:"required"`
	jwt.StandardClaims
}

// ParseToken 解析JWT，但不验证签名
func ParseToken(tokenStr string) (claims MyClaims, err error) {
	if _, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &claims); err != nil {
		return claims, errorx.ErrInvalidToken.WithMessage("invalid token")
	}
	return claims, nil
}

// ValidToken 校验JWT
func ValidToken(tokenStr string, secret string) (MyClaims, error) {
	// 解析token
	var claims MyClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	if err != nil {
		// token 过期校验
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return claims, errorx.ErrExpiredToken.WithMessage("token is expired")
			}
		}
		logz.WarnNoCtx(fmt.Sprintf("invalid token %s", tokenStr))
		return claims, errorx.ErrInvalidToken.WithError(err)
	}
	// 校验token是否有效
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return *claims, nil
	} else {
		logz.WarnNoCtx(fmt.Sprintf("invalid token %s", tokenStr))
		return *claims, errorx.ErrInvalidToken.WithMessage("invalid token")
	}
}

// RefreshToken 刷新JWT
func RefreshToken(claims MyClaims, secret string) (string, error) {
	expirationTime := time.Unix(claims.StandardClaims.ExpiresAt, 0)
	timeToExpire := expirationTime.Sub(time.Now())

	// 如果距离过期时间不足 10分钟，则刷新 token
	if timeToExpire < TokenRefreshDuration {
		return GenToken(claims, secret)
	}
	// 否则，返回空字符串表示 token 未过期，不需要刷新
	return "", nil
}

// GenToken 生成JWT
func GenToken(claims MyClaims, secret string) (string, error) {
	// 创建一个我们自己的声明
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
		Issuer:    "Venus SaaS WAF",                           // 签发人
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(secret))
}
