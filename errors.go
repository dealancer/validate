package validate

import (
	"fmt"
	"reflect"
)

// ErrorValidation occurs when validator does not validate.
type ErrorValidation struct {
	FieldName      string
	FieldValue     reflect.Value
	ValidatorType  ValidatorType
	ValidatorValue string
}

func (e ErrorValidation) Error() string {
	validator := string(e.ValidatorType)
	if len(e.ValidatorValue) > 0 {
		validator += "=" + e.ValidatorValue
	}

	if len(e.FieldName) > 0 {
		return fmt.Sprintf("Validation error in field \"%v\" of type \"%v\" using validator \"%v\"", e.FieldName, e.FieldValue.Type(), validator)
	}

	return fmt.Sprintf("Validation error in value of type \"%v\" using validator \"%v\"", e.FieldValue.Type(), validator)
}
