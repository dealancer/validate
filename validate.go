package validate

import (
	"errors"
	"reflect"
	"strings"
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

// validateStruct iterates over struct fields
func validateStruct(value reflect.Value) error {
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		if err := validateField(value.Field(i), typ.Field(i), false); err != nil {
			return err
		}
	}

	return nil
}

// validateField validates a struct field
func validateField(value reflect.Value, field reflect.StructField, isChild bool) error {
	kind := value.Kind()
	tag := field.Tag

	// Perform validators
	valMap := parseValidateTag(tag)
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

// parseValidateTag parses validate tag into hash map of validators
func parseValidateTag(tag reflect.StructTag) map[string]string {
	valMap := make(map[string]string)
	entries := strings.Split(tag.Get(masterTag), ",")
	for _, e := range entries {
		parts := strings.Split(e, "=")

		if len(parts) == 2 {
			n := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])

			if n != "" {
				valMap[n] = v
			}
		}
	}

	return valMap
}

// getValidatorType returns validator type
func getValidatorType(valType string, child bool) string {
	var valChildMap = map[string]string{
		valTypeChildMin:   valTypeMin,
		valTypeChildMax:   valTypeMax,
		valTypeChildEmpty: valTypeEmpty,
		valTypeChildNil:   valTypeNil,
		valTypeChildOneOf: valTypeOneOf,
	}

	if child {
		if valType, ok := valChildMap[valType]; ok {
			return valType
		}
		return ""
	}

	return valType
}

// parseTokens splits validator value into tokens
func parseTokens(str string) []interface{} {
	if strings.TrimSpace(str) == "" {
		return nil
	}

	tokenStrings := strings.Split(str, "|")
	tokens := make([]interface{}, len(tokenStrings))
	for i := range tokenStrings {
		tokens[i] = strings.TrimSpace(tokenStrings[i])
	}

	return tokens
}

// isOneOf check if a token is one of tokens
func isOneOf(token interface{}, tokens []interface{}) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}

	return false
}
