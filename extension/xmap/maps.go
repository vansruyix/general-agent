// Package xmap
// @author: liulun
// @date: 2024/4/16
// @note: 自定义map列表
package xmap

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONMaps 是一个自定义类型，表示一个包含多个 map[string]interface{} 的切片
type JSONMaps []map[string]interface{}

// Value 实现 driver.Valuer 接口，将 JSONMaps 转换为数据库支持的值
func (j JSONMaps) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口，将数据库中的值转换为 JSONMaps 类型
func (j *JSONMaps) Scan(value interface{}) error {
	// 检查 value 是否为空
	if value == nil {
		*j = nil
		return nil
	}

	// 将数据库中的值转换为 []byte
	b, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source")
	}

	// 解码 JSON 数据并存储到 j 中
	return json.Unmarshal(b, j)
}
