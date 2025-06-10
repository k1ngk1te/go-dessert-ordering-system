package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	return &Validator{validator: v}
}

func (v *Validator) ValidateStruct(s interface{}) map[string]string {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	validatorErrors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := strings.ToLower(err.Field())
		validatorErrors[fieldName] = fmt.Sprintf("invalid %s: %s", fieldName, err.Tag())
		// You can customize error messages more precisely here.
		// For example:
		// switch err.Tag() {
		// case "required":
		//     validationErrors[fieldName] = fmt.Sprintf("%s is required", fieldName)
		// case "email":
		//     validationErrors[fieldName] = fmt.Sprintf("%s must be a valid email address", fieldName)
		// case "min":
		//     validationErrors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param())
		// default:
		//     validationErrors[fieldName] = fmt.Sprintf("invalid %s: %s", fieldName, err.Tag())
		// }
	}

	return validatorErrors
}
