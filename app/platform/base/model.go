package base

import "github.com/go-resty/resty/v2"

type FileReq struct {
	ContentType string
	Body        []byte // 文件内容
}
type FileReqCommon interface {
	GetContentType() string
	GetBody() []byte
}

func (r *FileReq) GetContentType() string {
	return r.ContentType
}
func (r *FileReq) GetBody() []byte {
	return r.Body
}

// FileResp 文件响应.接收resty的文件响应
type FileResp struct {
	FileName  string
	RestyResp *resty.Response
	Body      []byte // 文件内容。从Response.RawBody()获取
}

type FileRespCommon interface {
	GetRestyResp() *resty.Response
	SetRestyResp(resp *resty.Response)
	GetBody() []byte
	SetBody(body []byte)
}

func (r *FileResp) GetRestyResp() *resty.Response {
	return r.RestyResp
}
func (r *FileResp) SetRestyResp(resp *resty.Response) {
	r.RestyResp = resp
}
func (r *FileResp) GetBody() []byte {
	return r.Body
}
func (r *FileResp) SetBody(body []byte) {
	r.Body = body
}
