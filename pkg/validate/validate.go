package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-rock/rock/binding"

	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func ValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errs := err.(validator.ValidationErrors)
			for _, e := range errs {
				transtr := e.Translate(trans)
				f := strings.ToLower(e.StructField())
				res.Errors[f] = transtr
			}
		default:
			res.Errors["error"] = err.Error()
		}
	}
	return res
}

func InitBinding() {
	zhs := zh.New()
	uni = ut.New(zhs, zhs)

	trans, _ = uni.GetTranslator("zh")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniq", ValidateUniq)
		v.RegisterTranslation("uniq", trans, func(ut ut.Translator) error {
			return ut.Add("uniq", "{0}已经存在!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("uniq", fe.Field())
			return t
		})
		zh_translations.RegisterDefaultTranslations(v, trans)
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("comment")
		})
	}

}
