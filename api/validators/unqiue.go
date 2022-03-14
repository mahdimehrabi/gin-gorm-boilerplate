package validators

import (
	"boilerplate/infrastructure"

	"github.com/go-playground/validator/v10"
)

type UniqueValidator struct {
	logger   infrastructure.Logger
	database infrastructure.Database
}

func NewUniqueValidator(logger infrastructure.Logger, database infrastructure.Database) UniqueValidator {
	return UniqueValidator{
		logger:   logger,
		database: database,
	}
}

func (uv UniqueValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		uv.logger.Zap.Warn("It works!!! :adult:")
		return false
	}
}
