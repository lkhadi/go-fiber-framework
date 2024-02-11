package requests

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	expectedFormat := "2006-01-02"
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)

	if !match {
		return false
	}

	_, err := time.Parse(expectedFormat, dateStr)
	return err == nil
}

func ValDatetime(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	expectedFormat := "2006-01-02 15:04:05"
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, dateStr)

	if !match {
		return false
	}

	_, err := time.Parse(expectedFormat, dateStr)
	return err == nil
}

func ValIn(fl validator.FieldLevel) bool {
	allowedValues := fl.Param()
	allowedSlice := strings.Split(allowedValues, "&")
	fieldValue := fl.Field().String()

	for _, allowedValue := range allowedSlice {
		if fieldValue == allowedValue {
			return true
		}
	}

	return false
}

func ValRequiredIf(fl validator.FieldLevel) bool {
	tagValue := fl.Param()
	tagParts := strings.Split(tagValue, "&")
	fieldValue := fl.Field().String()

	if len(tagParts) != 3 {
		return false
	}

	targetFieldName := tagParts[0]
	targetField := fl.Parent().FieldByName(targetFieldName)
	if !targetField.IsValid() {
		return false
	}

	if tagParts[1] == "=" {
		if (tagParts[2] == targetField.String() && fieldValue != "") || tagParts[2] != targetField.String() {
			return true
		}
		return false
	} else if tagParts[1] == "!=" {
		if (tagParts[2] != targetField.String() && fieldValue != "") || (tagParts[2] == targetField.String()) {
			return true
		}

		return false
	} else {
		return false
	}
}

func ValFileType(fl validator.FieldLevel) bool {
	allowedValues := fl.Param()
	allowedSlice := strings.Split(allowedValues, "&")
	fileHeader, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return true
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		fmt.Println(err)
		return false
	}

	contentType := http.DetectContentType(buffer)
	for _, t := range allowedSlice {
		if contentType == t {
			return true
		}
	}

	return false
}
