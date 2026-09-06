// Package valid
// @author: fengyi
// @date: 2024/8/28
// @note:
package valid

import (
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

func MaxBytes(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	maxBytes, err := parseSize(fl.Param())
	if err != nil {
		return false
	}
	return int64(len(str)) <= maxBytes
}

func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToLower(sizeStr)
	var multiplier int64 = 1
	if strings.HasSuffix(sizeStr, "k") {
		sizeStr = strings.TrimSuffix(sizeStr, "k")
		multiplier = 1024
	} else if strings.HasSuffix(sizeStr, "m") {
		sizeStr = strings.TrimSuffix(sizeStr, "m")
		multiplier = 1024 * 1024
	}
	// else if strings.HasSuffix(sizeStr, "g") {
	//	sizeStr = strings.TrimSuffix(sizeStr, "g")
	//	multiplier = 1024 * 1024 * 1024
	//}
	size, err := strconv.ParseInt(sizeStr, 10, 64)

	if err != nil {
		return 0, err
	}
	return size * multiplier, nil
}
