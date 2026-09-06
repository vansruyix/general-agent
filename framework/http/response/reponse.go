package response

import (
	"fmt"
	"general-agent/extension/errorx"
)

const ResponseDefaultMsg = "ok"
const ResponseDefaultErrorMsg = "fail"

type Response struct {
	RequestId   string `json:"requestId,omitempty"` //请求ID
	Code        int    `json:"code"`                //状态码
	Msg         string `json:"msg,omitempty"`       //信息
	Reason      string `json:"reason,omitempty"`    //错误信息
	Data        any    `json:"data,omitempty"`      //数据
	InternalErr error  `json:"-"`                   //内部错误
}

func (r Response) Error() error {
	if r.InternalErr != nil {
		if errorx.IsBizFault(r.InternalErr) {
			return r.InternalErr
		} else {
			return errorx.ErrDeviceRpcError.WithError(r.InternalErr)
		}
	}
	if r.Code != 0 {
		return errorx.ErrDeviceRpcError.WithReason(fmt.Sprintf("%s, code=%d, %s", r.Msg, r.Code, r.Reason))
	}
	return nil
}
