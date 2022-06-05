package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewUniqueValidator),
	fx.Provide(NewFkValidator),
	fx.Provide(NewTimestampValidator),
	fx.Provide(NewValidators),
)

type Validators struct {
	uv  UniqueValidator
	fkv FkValidator
	ts  TimestampValidator
}

func NewValidators(uv UniqueValidator, fkv FkValidator, ts TimestampValidator) Validators {
	return Validators{
		uv:  uv,
		fkv: fkv,
		ts:  ts,
	}
}

// Setup sets up middlewares
func (val Validators) Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueDB", val.uv.Handler())
		v.RegisterValidation("fkDB", val.fkv.Handler())
		v.RegisterValidation("timestamp", val.ts.Handler())
	}
}
