package customValidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	V                *validator.Validate
	operationTypeErr error
}

func NewCustomValidator() *CustomValidator {

	cv := &CustomValidator{V: validator.New()}

	cv.V.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := cv.V.RegisterValidation("operationType", cv.operationTypeValidate); err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) Validating(obj any) error {
	err := cv.V.Struct(obj)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			fieldErr := validationErrs[0]
			return cv.validationError(fieldErr.Field(), fieldErr.Tag())
		}
		return err
	}
	return nil
}

func (cv *CustomValidator) operationTypeValidate(fl validator.FieldLevel) bool {

	if fl.Field().Kind() != reflect.String {
		cv.operationTypeErr = fmt.Errorf("field %s must be a string", fl.FieldName())
		return false
	}

	fieldValue := fl.Field().String()

	if fieldValue != "DEPOSIT" && fieldValue != "WITHDRAW" {
		cv.operationTypeErr = fmt.Errorf(`field %s must be only DEPOSIT or WITHDRAW`, fl.FieldName())
		return false
	}
	return true
}

func (cv *CustomValidator) validationError(field string, tag string) error {
	switch tag {
	case "required":
		return fmt.Errorf("field %s is required", field)
	case "uuid":
		return fmt.Errorf("field %s must be a valid uuid", field)
	case "operationType":
		return cv.operationTypeErr
	case "gt=0":
		return fmt.Errorf("field %s must be more than 0", field)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}
