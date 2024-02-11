package reqUser

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

type Save struct {
	Nama            string                `validate:"required"`
	NIP             string                `validate:"required"`
	Role            string                `validate:"required,in=adminsystem&superadmin&operator&mekanik&foreman&superintendent"`
	Email           string                `validate:"required"`
	Password        string                `validate:"required"`
	Department      string                `validate:"required"`
	Section         string                `validate:"required"`
	Jabatan         string                `validate:"required"`
	Tanda_Tangan    *multipart.FileHeader `validate:"omitempty,filetype=image/jpg&image/png&image/jpeg"`
	UUID_Perusahaan string                `validate:"required_if=Role&!=&adminsystem"`
}

func (request *Save) Validate() (map[string]string, error) {
	validate := validator.New()
	validate.RegisterValidation("filetype", requests.ValFileType)
	validate.RegisterValidation("in", requests.ValIn)
	validate.RegisterValidation("required_if", requests.ValRequiredIf)
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
