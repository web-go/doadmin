package validate

import (
	"github.com/web-go/doadmin/app/models"

	"github.com/go-playground/validator/v10"
)

func ValidateUniq(fl validator.FieldLevel) bool {
	result := 0

	value := fl.Field().String()   // value
	column := fl.StructFieldName() // column name

	m := fl.Top()

	id := m.Elem().FieldByName("ID").Uint()

	m.Elem().FieldByName("ID").SetUint(0)
	model := m.Interface()

	models.DB.Model(model).Where(column+" = ? and id != ?", value, id).Count(&result)
	m.Elem().FieldByName("ID").SetUint(id)
	dup := result > 0
	return !dup
}
