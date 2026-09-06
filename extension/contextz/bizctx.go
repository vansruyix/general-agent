package contextz

import "context"

const (
	CTX_KEY string = "venus-ctx"

	// 调用方应用名
	CALLER string = "caller"
	// 当前租户ID
	USERNAME string = "username"
	// 当前员工ID
	ROLES string = "roles"
	// 租户ID
	TENANT_ID string = "X-CUSTOMER-ID"
	// 用户ID
	USER_ID string = "X-USER-ID"
	// XDR
	XDR_TENANT_ID string = "Userregionid"
	XDR_USER_ID   string = "Userid"
)

type BizContext map[string]string

func (b BizContext) IsEmpty() bool {
	return len(b) == 0
}

func WithBizCtx(ctx context.Context, bizCtx BizContext) context.Context {
	return context.WithValue(ctx, CTX_KEY, bizCtx)
}

func GetBizCtx(ctx context.Context) BizContext {
	bc := ctx.Value(CTX_KEY)
	if v, ok := bc.(BizContext); ok {
		return v
	}
	return make(BizContext)
}

func GetCaller(ctx context.Context) string {
	return GetBizCtx(ctx)[CALLER]
}

func GetUsername(ctx context.Context) string {
	return GetBizCtx(ctx)[USERNAME]
}

func GetRoles(ctx context.Context) string {
	return GetBizCtx(ctx)[ROLES]
}

func GetTenantID(ctx context.Context) string {
	return GetBizCtx(ctx)[TENANT_ID]
}

func GetUserID(ctx context.Context) string {
	return GetBizCtx(ctx)[USER_ID]
}

func GetXdrTenantID(ctx context.Context) string {
	return GetBizCtx(ctx)[XDR_TENANT_ID]
}

func GetXdrUserID(ctx context.Context) string {
	return GetBizCtx(ctx)[XDR_USER_ID]
}

func (bc BizContext) GetCaller() string {
	return bc[CALLER]
}

func (bc BizContext) GetUsername() string {
	return bc[USERNAME]
}

func (bc BizContext) GetRoles() string {
	return bc[ROLES]
}

func (bc BizContext) GetTenantID() string {
	return bc[TENANT_ID]
}

func (bc BizContext) GetUserID() string {
	return bc[USER_ID]
}

func (bc BizContext) GetXdrTenantID() string {
	return bc[XDR_TENANT_ID]
}

func (bc BizContext) GetXdrUserID() string {
	return bc[XDR_USER_ID]
}
