package xresty

import (
	"crypto/tls"
	"fmt"
	"general-agent/app/platform/base"
	"general-agent/config"
	"general-agent/extension/errorx"
	"general-agent/extension/logz"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	Caller *resty.Client
}

func NewHttpClient(schema *config.Config) *HttpClient {
	caller := resty.New().EnableTrace()
	// Debug模式
	caller.SetDebug(schema.Debug)
	caller.SetTimeout(60 * time.Second)
	// 忽略证书验证
	caller.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	return &HttpClient{Caller: caller}
}

func (c *HttpClient) Get(url string, headers map[string]string, params map[string]string, result interface{}) (*resty.Response, error) {
	request := c.Caller.R()
	if len(params) > 0 {
		request.SetQueryParams(params)
	}
	if len(headers) > 0 {
		request.SetHeaders(headers)
	}

	// 文件下载
	if fr, ok := result.(base.FileRespCommon); ok {
		request.SetDoNotParseResponse(true).SetDebug(false)
		resp, err := request.Get(url)
		if err != nil {
			logz.ErrorNoCtx("[RPC] call error", "url", resp.Request.URL, "err", err)
			return resp, errorx.ErrDefault.WithMessage("接口调用失败")
		}
		defer func() {
			if fr.GetRestyResp() != nil {
				_ = fr.GetRestyResp().RawBody().Close()
			}
		}()

		if resp.StatusCode() != http.StatusOK {
			if resp.RawBody() != nil {
				_ = resp.RawBody().Close()
			}
			return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode()))
		}

		bodyBytes, err := io.ReadAll(resp.RawBody())
		if err != nil {
			logz.ErrorNoCtx("[RPC] read body error", "url", resp.Request.URL, "err", err)
			return resp, errorx.ErrDefault.WithMessage("接口调用失败")
		}

		fr.SetRestyResp(resp)
		fr.SetBody(bodyBytes)
		return resp, nil
	}

	resp, err := request.SetResult(result).Get(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] call error", "url", resp.Request.URL, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK {
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(resp.Body()))
	}
	return resp, nil
}

func (c *HttpClient) Post(url string, headers map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Caller.R().
		SetBody(body).
		SetHeaders(headers).
		SetResult(result).
		Post(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] call error", "url", resp.Request.URL, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(resp.Body()))
	}
	return resp, nil
}

func (c *HttpClient) Put(url string, headers map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Caller.R().
		SetBody(body).
		SetHeaders(headers).
		SetResult(result).
		Put(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] call error", "url", resp.Request.URL, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK {
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(resp.Body()))
	}
	return resp, nil
}

func (c *HttpClient) Delete(url string, headers map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Caller.R().
		SetBody(body).
		SetHeaders(headers).
		SetResult(result).
		Delete(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] call error", "url", resp.Request.URL, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(resp.Body()))
	}
	return resp, nil
}

func (c *HttpClient) Patch(url string, headers map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Caller.R().
		SetBody(body).
		SetHeaders(headers).
		SetResult(result).
		Patch(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] call error", "url", resp.Request.URL, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK {
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(resp.Body()))
	}
	return resp, nil
}

func HandleResponse(response *resty.Response, err error, respErr error) error {
	if err != nil {
		logz.ErrorNoCtx("[RPC] call error", "url", response.Request.URL, "err", err)
		return err
	}
	if response.StatusCode() != http.StatusOK {
		logz.ErrorNoCtx("[RPC] status error", "code", response.StatusCode(), "response", string(response.Body()))
		return fmt.Errorf("[RPC] status error %d", response.StatusCode())
	}
	if respErr != nil {
		logz.ErrorNoCtx("[RPC] response error", "url", response.Request.URL, "err", respErr)
		return respErr
	}
	return nil
}

// Stream 执行POST流式请求（SSE），设置SetDoNotParseResponse(true)不自动解析响应
// 返回原始HTTP响应，由调用方通过resp.RawBody()逐行读取事件流
// 用于AI对话、Workflow运行等需要实时输出的场景
func (c *HttpClient) Stream(url string, headers map[string]string, body interface{}) (*resty.Response, error) {
	request := c.Caller.R().
		SetBody(body).
		SetHeaders(headers).
		SetDoNotParseResponse(true)

	resp, err := request.Post(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] stream call error", "url", url, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.RawBody())
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(bodyBytes))
	}
	return resp, nil
}

// StreamGet 执行GET流式请求（SSE），设置SetDoNotParseResponse(true)不自动解析响应
// 返回原始HTTP响应，由调用方通过resp.RawBody()逐行读取事件流
// 用于需要实时输出的GET请求场景
func (c *HttpClient) StreamGet(url string, headers map[string]string, params map[string]string) (*resty.Response, error) {
	request := c.Caller.R().
		SetQueryParams(params).
		SetHeaders(headers).
		SetDoNotParseResponse(true)

	resp, err := request.Get(url)
	if err != nil {
		logz.ErrorNoCtx("[RPC] stream get error", "url", url, "err", err)
		return resp, errorx.ErrDefault.WithMessage("接口调用失败")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, errorx.ErrUnauthorized
	}
	if resp.StatusCode() != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.RawBody())
		return resp, errorx.ErrRpcError.WithMessage(strconv.Itoa(resp.StatusCode())).WithReason(string(bodyBytes))
	}
	return resp, nil
}
