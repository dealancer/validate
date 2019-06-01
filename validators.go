package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func validateMin(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if min, err := time.ParseDuration(validator); err == nil && time.Duration(value.Int()) < min {
				return errors.New(fmt.Sprint(name, " must not be less than ", min))
			}
		} else {
			if min, err := strconv.ParseInt(validator, 10, 64); err == nil && value.Int() < min {
				return errors.New(fmt.Sprint(name, " must not be less than ", min))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if min, err := strconv.ParseUint(validator, 10, 64); err == nil && value.Uint() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
	case reflect.Float32, reflect.Float64:
		if min, err := strconv.ParseFloat(validator, 64); err == nil && value.Float() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
	case reflect.String:
		if min, err := strconv.Atoi(validator); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " characters"))
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		if min, err := strconv.Atoi(validator); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " elements"))
		}
	}

	return nil
}

func validateMax(value reflect.Value, name string, validator string) error {
	kind := value.Kind()
	typ := value.Type()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if max, err := time.ParseDuration(validator); err == nil && time.Duration(value.Int()) > max {
				return errors.New(fmt.Sprint(name, " must not be greater than ", max))
			}
		} else {
			if max, err := strconv.ParseInt(validator, 10, 64); err == nil && value.Int() > max {
				return errors.New(fmt.Sprint(name, " must not be greater than ", max))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if max, err := strconv.ParseUint(validator, 10, 64); err == nil && value.Uint() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.Float32, reflect.Float64:
		if max, err := strconv.ParseFloat(validator, 64); err == nil && value.Float() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.String:
		if max, err := strconv.Atoi(validator); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " characters"))
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		if max, err := strconv.Atoi(validator); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " elements"))
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
				if !isOneOf(time.Duration(value.Int()), tokens) {
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
				if !isOneOf(value.Int(), tokens) {
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
			if !isOneOf(value.Uint(), tokens) {
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
			if !isOneOf(value.Float(), tokens) {
				return errors.New(fmt.Sprint(name, " must be one of ", validator))
			}
		}
	case reflect.String:
		if tokens := parseTokens(validator); len(tokens) > 0 {
			if !isOneOf(value.String(), tokens) {
				return errors.New(fmt.Sprint(name, " must be one of ", validator))
			}
		}
	}

	return nil
}
