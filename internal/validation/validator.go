package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := &Validator{
		validate: validator.New(),
	}
	
	// Register custom validation functions
	v.registerCustomValidations()
	
	return v
}

func (v *Validator) registerCustomValidations() {
	// Example custom validation: horse gender
	v.validate.RegisterValidation("horse_gender", func(fl validator.FieldLevel) bool {
		gender := fl.Field().String()
		return gender == "MARE" || gender == "STALLION" || gender == "GELDING"
	})

	// Example custom validation: positive number
	v.validate.RegisterValidation("positive", func(fl validator.FieldLevel) bool {
		val := fl.Field().Float()
		return val > 0
	})
}

func (v *Validator) Struct(s interface{}) error {
	if s == nil {
		return fmt.Errorf("cannot validate nil struct")
	}

	err := v.validate.Struct(s)
	if err != nil {
		return v.processValidationErrors(err)
	}
	return nil
}

func (v *Validator) processValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var errorMessages []string
	for _, e := range validationErrors {
		errorMessages = append(errorMessages, v.formatFieldError(e))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
}

func (v *Validator) formatFieldError(e validator.FieldError) string {
	fieldName := e.Field()
	fieldValue := e.Value()
	errorType := e.Tag()

	switch errorType {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "min":
		return fmt.Sprintf("%s must be greater than or equal to %s", fieldName, e.Param())
	case "max":
		return fmt.Sprintf("%s must be less than or equal to %s", fieldName, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fieldName, e.Param())
	default:
		return fmt.Sprintf("%s failed %s validation (value: %v)", fieldName, errorType, fieldValue)
	}
}

// ValidateStruct is a convenience function for global validation
func ValidateStruct(s interface{}) error {
	return New().Struct(s)
}
