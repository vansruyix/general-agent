// Package bean
// @author: liulun
// @date: 2024/7/2
// @note: 字符串数组
package bean

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StrArray []string

func (f *StrArray) Scan(value interface{}) error {
	if value == nil {
		*f = StrArray{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source")
	}
	return json.Unmarshal(bytes, f)
}

// Value 实现 driver.Valuer 接口
func (f StrArray) Value() (driver.Value, error) {
	if len(f) == 0 {
		return "[]", nil
	}
	jsonValue, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return string(jsonValue), nil
}
