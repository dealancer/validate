package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// ValidatorType is used for validator type definitions.
type ValidatorType string

// Following validators are available.
const (
	// ValidatorGt (greater than) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"gte=0"`
	ValidatorGt ValidatorType = "gt"

	// ValidatorLt (less than) compares a numeric value of a number or compares a count of elements in a string, a map, a slice, or an array.
	// E.g. `validate:"lte=10"`
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
type validatorFunc func(value reflect.Value, name string, validator string) error

func getValidatorTypeMap() map[ValidatorType]validatorFunc {
	return map[ValidatorType]validatorFunc{
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

func validateGt(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err == nil && time.Duration(value.Int()) <= token {
				return errors.New(fmt.Sprint(name, " must be greater than ", token))
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err == nil && value.Int() <= token {
				return errors.New(fmt.Sprint(name, " must be greater than ", token))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err == nil && value.Uint() <= token {
			return errors.New(fmt.Sprint(name, " must be greater than ", token))
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err == nil && value.Float() <= token {
			return errors.New(fmt.Sprint(name, " must be greater than ", token))
		}
	case reflect.String:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() <= token {
			return errors.New(fmt.Sprint(name, " must contain greater than ", token, " characters"))
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() <= token {
			return errors.New(fmt.Sprint(name, " must contain greater than ", token, " elements"))
		}
	}

	return nil
}

func validateLt(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err == nil && time.Duration(value.Int()) >= token {
				return errors.New(fmt.Sprint(name, " must be less than ", token))
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err == nil && value.Int() >= token {
				return errors.New(fmt.Sprint(name, " must be less than ", token))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err == nil && value.Uint() >= token {
			return errors.New(fmt.Sprint(name, " must be less than ", token))
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err == nil && value.Float() >= token {
			return errors.New(fmt.Sprint(name, " must be less than ", token))
		}
	case reflect.String:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() >= token {
			return errors.New(fmt.Sprint(name, " must contain less than ", token, " characters"))
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() >= token {
			return errors.New(fmt.Sprint(name, " must contain less than ", token, " elements"))
		}
	}

	return nil
}

func validateGte(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err == nil && time.Duration(value.Int()) < token {
				return errors.New(fmt.Sprint(name, " must be greater than or equal to ", token))
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err == nil && value.Int() < token {
				return errors.New(fmt.Sprint(name, " must be greater than or equal to ", token))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err == nil && value.Uint() < token {
			return errors.New(fmt.Sprint(name, " must be greater than or equal to ", token))
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err == nil && value.Float() < token {
			return errors.New(fmt.Sprint(name, " must be greater than or equal to ", token))
		}
	case reflect.String:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() < token {
			return errors.New(fmt.Sprint(name, " must contain greater than or equal to ", token, " characters"))
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() < token {
			return errors.New(fmt.Sprint(name, " must contain greater than or equal to ", token, " elements"))
		}
	}

	return nil
}

func validateLte(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if token, err := time.ParseDuration(validator); err == nil && time.Duration(value.Int()) > token {
				return errors.New(fmt.Sprint(name, " must be less than or equal to ", token))
			}
		} else {
			if token, err := strconv.ParseInt(validator, 10, 64); err == nil && value.Int() > token {
				return errors.New(fmt.Sprint(name, " must be less than or equal to ", token))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if token, err := strconv.ParseUint(validator, 10, 64); err == nil && value.Uint() > token {
			return errors.New(fmt.Sprint(name, " must be less than or equal to ", token))
		}
	case reflect.Float32, reflect.Float64:
		if token, err := strconv.ParseFloat(validator, 64); err == nil && value.Float() > token {
			return errors.New(fmt.Sprint(name, " must be less than or equal to ", token))
		}
	case reflect.String:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() > token {
			return errors.New(fmt.Sprint(name, " must contain less than or equal to ", token, " characters"))
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		if token, err := strconv.Atoi(validator); err == nil && value.Len() > token {
			return errors.New(fmt.Sprint(name, " must contain less than or equal to ", token, " elements"))
		}
	}

	return nil
}

func validateEmpty(value reflect.Value, name string, validator string) error {
	kind := value.Kind()

	switch kind {
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if isEmpty, err := strconv.ParseBool(validator); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
	}

	return nil
}

func validateNil(value reflect.Value, name string, validator string) error {
	kind := value.Kind()

	switch kind {
	case reflect.Ptr:
		if isNil, err := strconv.ParseBool(validator); err == nil {
			if isNil && !value.IsNil() {
				return errors.New(fmt.Sprint(name, " must be nil"))
			} else if !isNil && value.IsNil() {
				return errors.New(fmt.Sprint(name, " must not be nil"))
			}
		}
	}

	return nil
}

func validateOneOf(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if tokens := parseTokens(validator); len(tokens) > 0 {
				for i, token := range tokens {
					tokens[i] = nil
					if token, err := time.ParseDuration(token.(string)); err == nil {
						tokens[i] = token
					}
				}
				if !tokenOneOf(time.Duration(value.Int()), tokens) {
					return errors.New(fmt.Sprint(name, " must be one of ", validator))
				}
			}
		} else {
			if tokens := parseTokens(validator); len(tokens) > 0 {
				for i, token := range tokens {
					tokens[i] = nil
					if token, err := strconv.ParseInt(token.(string), 10, 64); err == nil {
						tokens[i] = token
					}
				}
				if !tokenOneOf(value.Int(), tokens) {
					return errors.New(fmt.Sprint(name, " must be one of ", validator))
				}
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if tokens := parseTokens(validator); len(tokens) > 0 {
			for i, token := range tokens {
				tokens[i] = nil
				if token, err := strconv.ParseUint(token.(string), 10, 64); err == nil {
					tokens[i] = token
				}
			}
			if !tokenOneOf(value.Uint(), tokens) {
				return errors.New(fmt.Sprint(name, " must be one of ", validator))
			}
		}
	case reflect.Float32, reflect.Float64:
		if tokens := parseTokens(validator); len(tokens) > 0 {
			for i, token := range tokens {
				tokens[i] = nil
				if token, err := strconv.ParseFloat(token.(string), 64); err == nil {
					tokens[i] = token
				}
			}
			if !tokenOneOf(value.Float(), tokens) {
				return errors.New(fmt.Sprint(name, " must be one of ", validator))
			}
		}
	case reflect.String:
		if tokens := parseTokens(validator); len(tokens) > 0 {
			if !tokenOneOf(value.String(), tokens) {
				return errors.New(fmt.Sprint(name, " must be one of ", validator))
			}
		}
	}

	return nil
}

func validateFormat(value reflect.Value, name string, validator string) error {
	kind := value.Kind()

	switch kind {
	case reflect.String:
		formatTypeMap := getFormatTypeMap()
		if formatFunc, ok := formatTypeMap[FormatType(validator)]; ok {
			if !formatFunc(value.String()) {
				return errors.New(fmt.Sprint(name, " is not valid ", validator))
			}
		}
	}

	return nil
}
