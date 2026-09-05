// Package errs 提供统一的业务错误码、响应封装与辅助函数。
// 所有 controller 通过 errs.OK / errs.Fail 返回格式化 JSON 响应，
// 业务层通过预定义的 BizError（如 ErrNotFound）传递错误语义。
package errs
