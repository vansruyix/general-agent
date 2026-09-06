// Package valid
// @author: liulun
// @date: 2024/8/21
// @note: 包含key
package valid

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func SkipZero(fl validator.FieldLevel) bool {
	fmt.Println(fl.FieldName())
	param := fl.Param()
	parts := strings.Split(param, ";")
	//如果为zero就跳过验证
	if fl.Field().IsZero() {
		return true
	}
	fieldValue := fl.Field().Interface()
	for i := 1; i < len(parts); i++ {
		rule := parts[i]
		if validator.New().Var(fieldValue, rule) != nil {
			return false
		}
	}
	return true
}
