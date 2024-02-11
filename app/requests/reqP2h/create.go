package reqP2h

import (
	"fmt"
	"p2h-api/app/requests"
	"strings"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type Save struct {
	UUID  string        `validate:"required"`
	Jam   string        `validate:"required,datetime"`
	Hm_Km string        `validate:"required"`
	Item  []ItemP2HSave `validate:"required,dive"`
}

type ItemP2HSave struct {
	UUID    string `validate:"required"`
	Kondisi string `validate:"required,in=Baik/Normal&Rusak/Tidak Normal"`
	Catatan string
}

func (request *Save) Validate() (map[string]string, error) {
	validate := validator.New()
	validate.RegisterValidation("datetime", requests.ValDatetime)
	validate.RegisterValidation("in", requests.ValIn)

	locale := id.New()
	uni := ut.New(locale, locale)
	trans, _ := uni.GetTranslator("id_ID")

	if err := idTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}

	err := validate.Struct(request)

	if err != nil {
		errors := make(map[string]string)
		validationErrors := err.(validator.ValidationErrors)
		validationErrors.Translate(trans)

		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.StructField()
			translatedError := fieldErr.Translate(trans)
			errors[strings.ToLower(fieldName)] = translatedError
			if fieldErr.Tag() == "datetime" {
				m := fmt.Sprintf("%s harus sesuai dengan format YYYY-MM-DD HH:ii:ss", fieldName)
				errors[strings.ToLower(fieldName)] = m
			}
		}
		return errors, fmt.Errorf("validation error")
	}

	return nil, nil
}
