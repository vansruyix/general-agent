// Package xjson
// @author: liulun
// @date: 2024/3/21
// @note: json对象
package xjson

import (
	"encoding/json"
	"general-agent/extension/errorx"
	"strconv"
)

type JsonObject struct {
	data map[string]string
}

func MakeJsonObject(jsonStr string) JsonObject {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		panic(errorx.ErrDefault.WithError(err).WithMessage("Marshal error"))
	}
	jsonObject := JsonObject{data: m}
	return jsonObject
}
func (j JsonObject) GetInt(name string) (int, error) {
	value := j.data[name]
	return strconv.Atoi(value)
}
func (j JsonObject) GetString(name string) string {
	return j.data[name]
}
