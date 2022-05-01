package validators

import (
	"boilerplate/core/infrastructure"

	"github.com/go-playground/validator/v10"
)

type FkValidator struct {
	logger   infrastructure.Logger
	database infrastructure.Database
}

func NewFkValidator(logger infrastructure.Logger, database infrastructure.Database) FkValidator {
	return FkValidator{
		logger:   logger,
		database: database,
	}
}

//fk validator
//please send destionation table name as foreign key
func (uv FkValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field()
		table := fl.Param()
		var count int64
		uv.database.DB.Table(table).Where("id=?", value.Uint()).Count(&count)
		return count > 0
	}
}
