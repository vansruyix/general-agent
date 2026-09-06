// Package valid
// @author: fengyi
// @date: 2024/7/2
// @note:
package valid

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Module("validator",
	fx.Invoke(registerValidation),
)

func registerValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validate_if", ValidateIf)
		if err != nil {
			panic(err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("internal_filename", InternalFileName)
		if err != nil {
			panic(err)
		}
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("skip_zero", SkipZero)
		if err != nil {
			panic(err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("maxbytes", MaxBytes)
		if err != nil {
			panic(err)
		}
	}
}
