package server

import (
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

func validateRequest(dataStruct interface{}) error {
	enUs := en.New()
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("label")
	})

	uni := ut.New(enUs)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	if err := validate.Struct(dataStruct); err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			return errors.New(validationErr.Translate(trans))
		}
	}
	return nil
}
