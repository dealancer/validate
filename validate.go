package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	masterTag = "validate"

	valMin     = "min"
	valMax     = "max"
	valIsEmpty = "is_empty"
	valIsNil   = "is_nil"

	valChildMin     = "child_min"
	valChildMax     = "child_max"
	valChildIsEmpty = "child_is_empty"
	valChildIsNil   = "child_is_nil"
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

	return errors.New(fmt.Sprint("not a struct or a struct pointer"))
}

func validateStruct(value reflect.Value) error {
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		if err := validateField(value.Field(i), typ.Field(i), false); err != nil {
			return err
		}
	}

	return nil
}

func validateField(value reflect.Value, field reflect.StructField, isChild bool) error {
	kind := value.Kind()
	typ := value.Type()
	name := field.Name
	tag := field.Tag
	valMap := parseVals(tag)

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if min, err := time.ParseDuration(getVal(valMap, valMin, isChild)); err == nil && time.Duration(value.Int()) < min {
				return errors.New(fmt.Sprint(name, " must not be less than ", min))
			}
		} else {
			if min, err := strconv.ParseInt(getVal(valMap, valMin, isChild), 10, 64); err == nil && value.Int() < min {
				return errors.New(fmt.Sprint(name, " must not be less than ", min))
			}
		}
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if max, err := time.ParseDuration(getVal(valMap, valMax, isChild)); err == nil && time.Duration(value.Int()) > max {
				return errors.New(fmt.Sprint(name, " must not be greater than ", max))
			}
		} else {
			if max, err := strconv.ParseInt(getVal(valMap, valMax, isChild), 10, 64); err == nil && value.Int() > max {
				return errors.New(fmt.Sprint(name, " must not be greater than ", max))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if min, err := strconv.ParseUint(getVal(valMap, valMin, isChild), 10, 64); err == nil && value.Uint() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseUint(getVal(valMap, valMax, isChild), 10, 64); err == nil && value.Uint() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.Float32, reflect.Float64:
		if min, err := strconv.ParseFloat(getVal(valMap, valMin, isChild), 64); err == nil && value.Float() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseFloat(getVal(valMap, valMax, isChild), 64); err == nil && value.Float() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.String:
		if isEmpty, err := strconv.ParseBool(getVal(valMap, valIsEmpty, isChild)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(getVal(valMap, valMin, isChild)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " characters"))
		}
		if max, err := strconv.Atoi(getVal(valMap, valMax, isChild)); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " characters"))
		}
	case reflect.Map:
		if isEmpty, err := strconv.ParseBool(getVal(valMap, valIsEmpty, isChild)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(getVal(valMap, valMin, isChild)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " elements"))
		}
		if max, err := strconv.Atoi(getVal(valMap, valMax, isChild)); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " elements"))
		}
	case reflect.Slice:
		if isEmpty, err := strconv.ParseBool(getVal(valMap, valIsEmpty, isChild)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(getVal(valMap, valMin, isChild)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " elements"))
		}
		if max, err := strconv.Atoi(getVal(valMap, valMax, isChild)); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " elements"))
		}
		if !isChild {
			for i := 0; i < value.Len(); i++ {
				if err := validateField(value.Index(i), field, true); err != nil {
					return err
				}
			}
		}
	case reflect.Ptr:
		if isNil, err := strconv.ParseBool(getVal(valMap, valIsNil, isChild)); err == nil {
			if isNil && !value.IsNil() {
				return errors.New(fmt.Sprint(name, " must be nil"))
			} else if !isNil && value.IsNil() {
				return errors.New(fmt.Sprint(name, " must not be nil"))
			}
		}
		if !value.IsNil() && !isChild {
			return validateField(value.Elem(), field, true)
		}
	}

	return nil
}

func parseVals(tag reflect.StructTag) map[string]string {
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

func getVal(valMap map[string]string, valName string, child bool) string {
	var valChildMap = map[string]string{
		valMin:     valChildMin,
		valMax:     valChildMax,
		valIsEmpty: valChildIsEmpty,
		valIsNil:   valChildIsNil,
	}

	if child {
		valName = valChildMap[valName]
	}

	if val, ok := valMap[valName]; ok {
		return val
	}

	return ""
}
