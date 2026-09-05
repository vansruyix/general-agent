package user

// CreateUserReq 是创建用户的请求体，绑定 JSON 请求。
type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	RoleID   int32  `json:"role_id"`
}

// UpdateUserReq 是更新用户的请求体，全部字段可选。
// RoleID 使用指针类型以区分"未传"与"传 0"。
type UpdateUserReq struct {
	Name   string `json:"name"`
	RoleID *int32 `json:"role_id"`
}

// ListUserReq 是分页查询用户的请求参数，从 query string 绑定。
type ListUserReq struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

// UserResp 是返回给前端的用户安全视图，不包含密码等敏感字段。
type UserResp struct {
	ID string `json:"id"`
	//用户名称
	Username   string `json:"username"`
	Name       string `json:"name"`
	RoleID     int32  `json:"role_id"`
	Status     int32  `json:"status"`
	CreateTime int64  `json:"create_time"`
}

// PageResp 是分页查询的通用响应结构。
type PageResp struct {
	List  []UserResp `json:"list"`
	Total int64      `json:"total"`
}

// toUserResp 将 User 实体转换为安全的 UserResp DTO，过滤敏感字段。
func toUserResp(u *User) UserResp {
	return UserResp{
		ID:         u.ID,
		Username:   u.Username,
		Name:       u.Name,
		RoleID:     u.RoleID,
		Status:     u.Status,
		CreateTime: u.CreateTime,
	}
}
