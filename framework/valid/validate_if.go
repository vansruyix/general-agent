// Package valid
// @author: fengyi
// @date: 2024/7/2
// @note:
package valid

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
)

func ValidateIf(fl validator.FieldLevel) bool {
	param := fl.Param()
	parts := strings.Split(param, ";")
	conditions := strings.Fields(parts[0])
	if len(conditions)%2 != 0 {
		return false
	}

	for i := 0; i < len(conditions); i += 2 {
		fieldName := conditions[i]
		expectedValue := conditions[i+1]
		fieldToCheck := fl.Parent().FieldByName(fieldName)
		if !fieldToCheck.IsValid() || !compareValues(fieldToCheck.Interface(), expectedValue) {
			// 如果条件满足就不进行后续验证
			return true
		}
	}
	fieldValue := fl.Field().Interface()
	for i := 1; i < len(parts); i++ {
		rule := parts[i]
		if validator.New().Var(fieldValue, rule) != nil {
			return false
		}
	}
	return true
}

func compareValues(actual interface{}, expected string) bool {
	switch v := actual.(type) {
	case string:
		return v == expected
	case int, int64, int32, int16, int8:
		vv, err := strconv.ParseInt(expected, 10, 64)
		if err != nil {
			return false
		}
		return reflect.ValueOf(v).Int() == vv
	case uint, uint64, uint32, uint16, uint8:
		vv, err := strconv.ParseUint(expected, 10, 64)
		if err != nil {
			return false
		}
		return reflect.ValueOf(v).Uint() == vv
	case float32, float64:
		vv, err := strconv.ParseFloat(expected, 64)
		if err != nil {
			return false
		}
		return reflect.ValueOf(v).Float() == vv
	case bool:
		vv, err := strconv.ParseBool(expected)
		if err != nil {
			return false
		}
		return v == vv
	default:
		return false
	}
}
