package xjson

import (
	"encoding/json"
	"general-agent/extension/errorx"
	"os"
)

func Marshal(obj any) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(errorx.ErrDefault.WithError(err).WithMessage("Marshal error"))
	}
	return string(bytes)
}
func MarshalIndent(obj any) string {
	bytes, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		panic(errorx.ErrDefault.WithError(err).WithMessage("Marshal error"))
	}
	return string(bytes)
}

func Unmarshal(data string, v any) {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		panic(errorx.ErrDefault.WithMessage("Unmarshal error"))
	}
}

func FileUnmarshal(file string, v any) error {
	_, err := os.Stat(file)
	if err != nil {
		return errorx.ErrDefault.WithError(err).WithMessage("file not exist")
	}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return errorx.ErrDefault.WithError(err).WithMessage("read file error")
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		return errorx.ErrDefault.WithError(err).WithMessage("file content format error")
	}
	return nil
}

func FileMarshal(file string, v any) error {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errorx.ErrDefault.WithError(err).WithMessage("file content format error")
	}
	if err := os.WriteFile(file, bytes, 0644); err != nil {
		return errorx.ErrDefault.WithError(err).WithMessage("write file error")
	}
	return nil
}

// Convert 将对象转换为指定结构
func Convert(obj any, target any) (any, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}
