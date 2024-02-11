package reqUser

import (
	"fmt"
	"p2h-api/app/requests"
	"strings"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type List struct {
	UUID_Perusahaan string
	Page            int `validate:"required,numeric"`
	Limit           int `validate:"required,numeric"`
	Search          string
	Sort            []requests.SortItem
}

func (request *List) Validate() (map[string]string, error) {
	validate := validator.New()
	validate.RegisterValidation("date", requests.ValDate)
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
			if fieldErr.Tag() == "date" {
				m := fmt.Sprintf("%s harus sesuai dengan format tanggal YYYY-MM-DD", fieldName)
				errors[strings.ToLower(fieldName)] = m
			}
		}
		return errors, fmt.Errorf("validation error")
	}

	return nil, nil
}
