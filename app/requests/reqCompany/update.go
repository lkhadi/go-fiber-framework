package reqCompany

import (
	"fmt"
	"mime/multipart"
	"p2h-api/app/requests"
	"strings"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type Update struct {
	Nama_Perusahaan string                `validate:"required"`
	Logo            *multipart.FileHeader `validate:"omitempty,filetype=image/jpg&image/png&image/jpeg"`
	UUID            string                `validate:"required"`
}

func (request *Update) Validate() (map[string]string, error) {
	validate := validator.New()
	validate.RegisterValidation("filetype", requests.ValFileType)
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
		}
		return errors, fmt.Errorf("validation error")
	}

	return nil, nil
}
