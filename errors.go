package validate

import (
	"fmt"
	"reflect"
)

// ErrorField is an error interface for field/value error.
type ErrorField interface {
	error
	FieldName() string
}

// errorField is a setter interface
type errorField interface {
	setFieldName(string)
}

// ErrorValidation occurs when validator does not validate.
type ErrorValidation struct {
	fieldName      string
	fieldValue     reflect.Value
	validatorType  ValidatorType
	validatorValue string
}

// FieldName gets a field name.
func (e ErrorValidation) FieldName() string {
	return e.fieldName
}

// setFieldName sets a field name.
func (e *ErrorValidation) setFieldName(fieldName string) {
	e.fieldName = fieldName
}

// Error returns an error.
func (e ErrorValidation) Error() string {
	validator := string(e.validatorType)
	if len(e.validatorValue) > 0 {
		validator += "=" + e.validatorValue
	}

	if len(e.fieldName) > 0 {
		return fmt.Sprintf("Validation error in field \"%v\" of type \"%v\" using validator \"%v\"", e.fieldName, e.fieldValue.Type(), validator)
	}

	return fmt.Sprintf("Validation error in value of type \"%v\" using validator \"%v\"", e.fieldValue.Type(), validator)
}

// ErrorSyntax occurs when there is a syntax error.
type ErrorSyntax struct {
	fieldName  string
	expression string
	near       string
	comment    string
}

// FieldName gets a field name.
func (e ErrorSyntax) FieldName() string {
	return e.fieldName
}

// setFieldName sets a field name.
func (e *ErrorSyntax) setFieldName(fieldName string) {
	e.fieldName = fieldName
}

// Error returns an error.
func (e ErrorSyntax) Error() string {
	if len(e.fieldName) > 0 {
		return fmt.Sprintf("Syntax error when validating field \"%v\", expression \"%v\" near \"%v\": %v", e.fieldName, e.expression, e.near, e.comment)
	}

	return fmt.Sprintf("Syntax error when validating value, expression \"%v\" near \"%v\": %v", e.expression, e.near, e.comment)
}

// Set field name
func setFieldName(err ErrorField, fieldName string) ErrorField {
	switch (err).(type) {
	case ErrorValidation:
		e := err.(ErrorValidation)
		var i interface{} = &e
		(i).(errorField).setFieldName(fieldName)
		return e
	case ErrorSyntax:
		e := err.(ErrorSyntax)
		var i interface{} = &e
		(i).(errorField).setFieldName(fieldName)
		return e
	}

	return err
}
