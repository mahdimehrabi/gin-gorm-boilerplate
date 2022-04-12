package validators

import (
	"boilerplate/core/infrastructure"
	"strings"

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

//unique validator
//please send table name and column name as parmater
//splited by & like this uniqueDB=users&email
//for update please define a ID prop and fill that in your controller
func (uv UniqueValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		//extract id for check unique of update
		var ID uint64 = 0
		IDField := fl.Parent().FieldByName("ID")
		if IDField.IsValid() && !IDField.IsZero() {
			ID = IDField.Uint()
		}

		value := fl.Field()
		str := strings.Split(fl.Param(), "&")
		table, column := str[0], str[1]
		var count int64
		uv.database.DB.Table(table).Where(column, value).Where("ID != ?", ID).Count(&count)
		return count < 1
	}
}
