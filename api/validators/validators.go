package validators

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewUniqueValidator),
	fx.Provide(NewValidators),
)

type Validators struct {
	uv UniqueValidator
}

func NewValidators(uv UniqueValidator) Validators {
	return Validators{
		uv: uv,
	}
}

// Setup sets up middlewares
func (val Validators) Setup() {
	fmt.Println("oooooooooooooooooooooooook")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		fmt.Println("nooooooooooooooooooooo")
		v.RegisterValidation("uniqueDB", val.uv.Handler())
	} else {
		fmt.Println("fuuuuuuuuuuuuuuuuuuuuuck")
	}
}
