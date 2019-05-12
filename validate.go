package validate

import (
	"errors"
	"reflect"
)

const (
	masterTag = "validate"

	valTypeMin   = "min"
	valTypeMax   = "max"
	valTypeEmpty = "empty"
	valTypeNil   = "nil"
	valTypeOneOf = "one_of"

	valTypeChildMin   = "child_min"
	valTypeChildMax   = "child_max"
	valTypeChildEmpty = "child_empty"
	valTypeChildNil   = "child_nil"
	valTypeChildOneOf = "child_one_of"
)

// Validate validates members of a struct
func Validate(element interface{}) error {
	value := reflect.ValueOf(element)

	if value.Kind() == reflect.Ptr {
		if value.Elem().Kind() == reflect.Struct {
			return validateStruct(value.Elem())
		}
	} else if value.Kind() == reflect.Struct {
		return validateStruct(value)
	}

	return errors.New("not a struct or a struct pointer")
}

// Iterate over struct fields
func validateStruct(value reflect.Value) error {
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		if err := validateField(value.Field(i), typ.Field(i), false); err != nil {
			return err
		}
	}

	return nil
}

// Validate struct field
func validateField(value reflect.Value, field reflect.StructField, isChild bool) error {
	kind := value.Kind()
	tag := field.Tag

	// Perform validators
	valMap := parseValidators(tag)
	for valType, validator := range valMap {
		valType := getValidatorType(valType, isChild)
		var err error

		switch valType {
		case valTypeMin:
			err = validateMin(value, field, validator)
		case valTypeMax:
			err = validateMax(value, field, validator)
		case valTypeEmpty:
			err = validateEmpty(value, field, validator)
		case valTypeNil:
			err = validateNil(value, field, validator)
		case valTypeOneOf:
			err = validateOneOf(value, field, validator)
		}

		if err != nil {
			return err
		}
	}

	// Dive one level deep into arrays and pointers
	switch kind {
	case reflect.Slice:
		if !isChild {
			for i := 0; i < value.Len(); i++ {
				if err := validateField(value.Index(i), field, true); err != nil {
					return err
				}
			}
		}
	case reflect.Ptr:
		if !value.IsNil() && !isChild {
			if err := validateField(value.Elem(), field, true); err != nil {
				return err
			}
		}
	}

	return nil
}
