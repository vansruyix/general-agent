package common

type ResultResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data any    `json:"data"`
}
