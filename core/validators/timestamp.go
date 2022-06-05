package validators

import (
	"boilerplate/core/infrastructure"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type TimestampValidator struct {
	logger   infrastructure.Logger
	database infrastructure.Database
}

func NewTimestampValidator(logger infrastructure.Logger, database infrastructure.Database) TimestampValidator {
	return TimestampValidator{
		logger:   logger,
		database: database,
	}
}

//timestamp validator
func (uv TimestampValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().Int()
		return len(strconv.Itoa(int(value))) == 10
	}
}
