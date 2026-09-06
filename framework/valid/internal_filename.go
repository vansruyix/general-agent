// Package valid
// @author: fengyi
// @date: 2024/7/10
// @note:
package valid

import (
	"general-agent/extension/xfile"

	"github.com/go-playground/validator/v10"
)

func InternalFileName(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	if fieldValue == "" {
		return false
	}
	return xfile.InternalFileNameRegexp.MatchString(fieldValue)
}
