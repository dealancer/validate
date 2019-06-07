package validate

import (
	"reflect"
	"strconv"
	"time"
)

// ValidatorType is used for validator type definitions.
type ValidatorType string

// Following validators are available.
const (
	// ValidatorEq (equals) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"eq=1"`
	ValidatorEq ValidatorType = "eq"

	// ValidatorNe (not equals) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"ne=0"`
	ValidatorNe ValidatorType = "ne"

	// ValidatorGt (greater than) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"gt=-1"`
	ValidatorGt ValidatorType = "gt"

	// ValidatorLt (less than) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"lt=11"`
	ValidatorLt ValidatorType = "lt"

	// ValidatorGte (greater than or equal to) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"gte=0"`
	ValidatorGte ValidatorType = "gte"

	// ValidatorLte (less than or equal to) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"lte=10"`
	ValidatorLte ValidatorType = "lte"

	// ValidatorEmpty checks if a string, a map, a slice, or an array is (not) empty.
	// E.g. `validate:"empty=false"`
	ValidatorEmpty ValidatorType = "empty"

	// ValidatorNil checks if a pointer is (not) nil.
	// E.g. `validate:"nil=false"`
	ValidatorNil ValidatorType = "nil"

	// ValidatorOneOf checks if a number or a string contains any of the given elements.
	// E.g. `validate:"one_of=1,2,3"`
	ValidatorOneOf ValidatorType = "one_of"

	// ValidatorFormat checks if a string of a given format.
	// E.g. `validate:"format=email"`
	ValidatorFormat ValidatorType = "format"
)

// validatorFunc is an interface for validator func
type validatorFunc func(value reflect.Value, validator string) ErrorField

func getValidatorTypeMap() map[ValidatorType]validatorFunc {
	return map[ValidatorType]validatorFunc{
		ValidatorEq:     validateEq,
		ValidatorNe:     validateNe,
		ValidatorGt:     validateGt,
		ValidatorLt:     validateLt,
		ValidatorGte:    validateGte,
		ValidatorLte:    validateLte,
		ValidatorEmpty:  validateEmpty,
		ValidatorNil:    validateNil,
		ValidatorOneOf:  validateOneOf,
		ValidatorFormat: validateFormat,
	}
}

type validator struct {
	Type  ValidatorType
	Value string
}

func validateEq(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorEq,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorEq),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err != nil {
				return errorSyntax
			} else if time.Duration(value.Int()) != token {
				return errorValidation
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err != nil {
				return errorSyntax
			} else if value.Int() != token {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err != nil {
			return errorSyntax
		} else if value.Uint() != token {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err != nil {
			return errorSyntax
		} else if value.Float() != token {
			return errorValidation
		}
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err != nil {
			return errorSyntax
		} else if value.Len() != token {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateNe(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorNe,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorNe),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err != nil {
				return errorSyntax
			} else if time.Duration(value.Int()) == token {
				return errorValidation
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err != nil {
				return errorSyntax
			} else if value.Int() == token {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err != nil {
			return errorSyntax
		} else if value.Uint() == token {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err != nil {
			return errorSyntax
		} else if value.Float() == token {
			return errorValidation
		}
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err != nil {
			return errorSyntax
		} else if value.Len() == token {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateGt(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorGt,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorGt),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err != nil {
				return errorSyntax
			} else if time.Duration(value.Int()) <= token {
				return errorValidation
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err != nil {
				return errorSyntax
			} else if value.Int() <= token {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err != nil {
			return errorSyntax
		} else if value.Uint() <= token {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err != nil {
			return errorSyntax
		} else if value.Float() <= token {
			return errorValidation
		}
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err != nil {
			return errorSyntax
		} else if value.Len() <= token {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateLt(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorLt,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorLt),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err != nil {
				return errorSyntax
			} else if time.Duration(value.Int()) >= token {
				return errorValidation
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err != nil {
				return errorSyntax
			} else if value.Int() >= token {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err != nil {
			return errorSyntax
		} else if value.Uint() >= token {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err != nil {
			return errorSyntax
		} else if value.Float() >= token {
			return errorValidation
		}
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err != nil {
			return errorSyntax
		} else if value.Len() >= token {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateGte(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorGte,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorGte),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err != nil {
				return errorSyntax
			} else if time.Duration(value.Int()) < token {
				return errorValidation
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err != nil {
				return errorSyntax
			} else if value.Int() < token {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err != nil {
			return errorSyntax
		} else if value.Uint() < token {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err != nil {
			return errorSyntax
		} else if value.Float() < token {
			return errorValidation
		}
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err != nil {
			return errorSyntax
		} else if value.Len() < token {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateLte(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorLte,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorLte),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err != nil {
				return errorSyntax
			} else if time.Duration(value.Int()) > token {
				return errorValidation
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err != nil {
				return errorSyntax
			} else if value.Int() > token {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err != nil {
			return errorSyntax
		} else if value.Uint() > token {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err != nil {
			return errorSyntax
		} else if value.Float() > token {
			return errorValidation
		}
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err != nil {
			return errorSyntax
		} else if value.Len() > token {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateEmpty(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorEmpty,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorEmpty),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if isEmpty, err := strconv.ParseBool(validator); err != nil {
			return errorSyntax
		} else if isEmpty && value.Len() > 0 {
			return errorValidation
		} else if !isEmpty && value.Len() == 0 {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateNil(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorNil,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorNil),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Ptr:
		if isNil, err := strconv.ParseBool(validator); err != nil {
			return errorSyntax
		} else if isNil && !value.IsNil() {
			return errorValidation
		} else if !isNil && value.IsNil() {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateOneOf(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()
	typ := value.Type()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorOneOf,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorOneOf),
		comment:    "could not parse or run",
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			var tokens []interface{}
			if tokens = parseTokens(validator); len(tokens) == 0 {
				return errorSyntax
			}
			for i, token := range tokens {
				tokens[i] = nil
				if token, err := time.ParseDuration(token.(string)); err != nil {
					return errorSyntax
				} else {
					tokens[i] = token
				}
			}
			if !tokenOneOf(time.Duration(value.Int()), tokens) {
				return errorValidation
			}
		} else {
			var tokens []interface{}
			if tokens = parseTokens(validator); len(tokens) == 0 {
				return errorSyntax
			}
			for i, token := range tokens {
				tokens[i] = nil
				if token, err := strconv.ParseInt(token.(string), 10, 64); err != nil {
					return errorSyntax
				} else {
					tokens[i] = token
				}
			}
			if !tokenOneOf(value.Int(), tokens) {
				return errorValidation
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var tokens []interface{}
		if tokens = parseTokens(validator); len(tokens) == 0 {
			return errorSyntax
		}
		for i, token := range tokens {
			tokens[i] = nil
			if token, err := strconv.ParseUint(token.(string), 10, 64); err != nil {
				return errorSyntax
			} else {
				tokens[i] = token
			}
		}
		if !tokenOneOf(value.Uint(), tokens) {
			return errorValidation
		}
	case reflect.Float32, reflect.Float64:
		var tokens []interface{}
		if tokens = parseTokens(validator); len(tokens) == 0 {
			return errorSyntax
		}
		for i, token := range tokens {
			tokens[i] = nil
			if token, err := strconv.ParseFloat(token.(string), 64); err != nil {
				return errorSyntax
			} else {
				tokens[i] = token
			}
		}
		if !tokenOneOf(value.Float(), tokens) {
			return errorValidation
		}
	case reflect.String:
		var tokens []interface{}
		if tokens = parseTokens(validator); len(tokens) == 0 {
			return errorSyntax
		}
		if !tokenOneOf(value.String(), tokens) {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}

func validateFormat(value reflect.Value, validator string) ErrorField {
	kind := value.Kind()

	errorValidation := ErrorValidation{
		fieldValue:     value,
		validatorType:  ValidatorFormat,
		validatorValue: validator,
	}

	errorSyntax := ErrorSyntax{
		expression: validator,
		near:       string(ValidatorFormat),
		comment:    "could not find format",
	}

	switch kind {
	case reflect.String:
		formatTypeMap := getFormatTypeMap()
		if formatFunc, ok := formatTypeMap[FormatType(validator)]; !ok {
			return errorSyntax
		} else if !formatFunc(value.String()) {
			return errorValidation
		}
	default:
		return errorSyntax
	}

	return nil
}
