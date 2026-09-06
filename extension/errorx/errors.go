package errorx

import (
	"fmt"
	"general-agent/extension/logz"
	"general-agent/extension/xstr"
	"strings"

	stderr "github.com/pkg/errors"
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

type Level int

type Error struct {
	// http response code, will be used in http response body.
	//
	// example: 404
	// if not set, default to DefaultSuccessCode
	Code int `json:"code,omitempty"` // 错误响应码

	// a human-readable reason of the error.
	//
	// example: "id not null"
	// required: true
	Reason string `json:"reason,omitempty"` // 错误的原因，用作开发排查问题

	// example: "the user does not exist"
	Message string `json:"message,omitempty"` // 通用错误信息，展示给用户

	err   error // 错误原始信息
	level Level // 错误级别
}

func (e Error) Error() string {
	buffer := strings.Builder{}
	if e.Message != "" {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(e.Message)
	}
	if e.Reason != "" {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(e.Reason)
	}
	if e.err != nil {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(e.err.Error())
	}
	return buffer.String()
}

func (e Error) Level() Level {
	return e.level
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Wrap(err error) error {
	e.err = err
	return WithStack(e)
}

func (e Error) WithWrap(err error) error {
	if err == nil {
		return WithStack(&e)
	}

	e.err = err
	return WithStack(&e)
}

func (e Error) WithError(err error) Error {
	if err == nil {
		return e
	}
	e.err = WithStack(err)
	return e
}

// WithReason reason
func (e Error) WithReason(reason string) Error {
	// reason 必须为英文。引文reason属于调试信息，不会经过国际化处理，避免语言为英文的时候，客户端在接口响应数据里看到中文字符
	// 如果业务没有多语言要求可以去掉
	if xstr.ContainsChinese(reason) {
		logz.WarnNoCtx("Unsatisfactory description!!!")
	}
	//if reason != "" {
	//	e.Reason = reason
	//}
	e.Reason = reason
	return e
}

func (e Error) GetReason() string {
	if e.Reason != "" {
		return e.Reason
	}
	if e.err != nil {
		return e.err.Error()

	}
	return ""
}

func (e Error) WithMessage(message string) Error {
	if message != "" {
		e.Message = message
	}
	return e
}
func (e Error) WithMessagef(message string, args ...any) Error {
	if message != "" {
		e.Message = fmt.Sprintf(message, args...)
	}
	return e
}

func (e Error) Cause() error {
	return e.err
}

func (e Error) As(value any) bool {
	switch v := value.(type) {
	case *Error:
		*v = e
		return true
	default:
		return false
	}
}

// New return a new Error object with error level
func New() Error {
	return Error{}
}

// NewError return a new Error object with error level
func NewError(code int, message string) Error {
	return NewErrorWithLevel(code, message, LevelError)
}

func NewErrorWithLevel(code int, message string, level Level) Error {
	return Error{Code: code, Message: message, level: level}
}

// IsBizFault 是否业务错误
func IsBizFault(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case Error, *Error:
		return true
	default:
		return false
	}
}

// WithStack mirrors the WithStack method of the stderr package.
// It adds a stack trace to the error if it does not already have one.
func WithStack(err error) error {
	if e, ok := err.(StackTracer); ok && len(e.StackTrace()) > 0 {
		return err
	}

	return stderr.WithStack(err)
}

func WithMessage(err error, msg string) error {
	return stderr.WithMessage(err, msg)
}

func WithMessageF(err error, format string, args ...interface{}) error {
	return stderr.WithMessagef(err, format, args...)
}

type StackTracer interface {
	StackTrace() stderr.StackTrace
}

type As interface {
	As(any) bool
}

type LevelContainer interface {
	Level() Level
}

func GetErrorLevel(err error, defaultLevel Level) Level {
	if err == nil {
		return defaultLevel
	}

	if e, ok := err.(LevelContainer); ok {
		return e.Level()
	}

	return defaultLevel
}
