package errs

// BizError 是携带业务错误码的结构化错误，实现 error 接口。
// Controller 层通过 errs.Fail 将 BizError 转为对应的 HTTP 200 + code/msg 响应。
type BizError struct {
	Code int
	Msg  string
}

// 预定义业务错误码。
// 错误码命名规则：1xxxx 为业务语义错误，50000 为内部错误。
var (
	ErrInternal   = &BizError{Code: 50000, Msg: "内部错误"}
	ErrNotFound   = &BizError{Code: 10001, Msg: "记录不存在"}
	ErrUserExists = &BizError{Code: 10002, Msg: "用户名已存在"}
)

// Error 实现 error 接口，返回错误消息。
func (e *BizError) Error() string {
	return e.Msg
}

// NewBizError 创建一个新的 BizError。
func NewBizError(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

func NewBizErrorMsg(msg string) *BizError {
	return &BizError{Code: 500, Msg: msg}
}
