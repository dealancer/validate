package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	tagMin     = "min"
	tagMax     = "max"
	tagIsEmpty = "is_empty"
	tagIsNil   = "is_nil"

	tagChildMin     = "child_min"
	tagChildMax     = "child_max"
	tagChildIsEmpty = "child_is_empty"
	tagChildIsNil   = "child_is_nil"
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

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if min, err := time.ParseDuration(getTag(tag, tagMin, isChild)); err == nil && time.Duration(value.Int()) < min {
				return errors.New(fmt.Sprint(name, " must not be less than ", min))
			}
		} else {
			if min, err := strconv.ParseInt(getTag(tag, tagMin, isChild), 10, 64); err == nil && value.Int() < min {
				return errors.New(fmt.Sprint(name, " must not be less than ", min))
			}
		}
		if typ == reflect.TypeOf((time.Duration)(0)) {
			if max, err := time.ParseDuration(getTag(tag, tagMax, isChild)); err == nil && time.Duration(value.Int()) > max {
				return errors.New(fmt.Sprint(name, " must not be greater than ", max))
			}
		} else {
			if max, err := strconv.ParseInt(getTag(tag, tagMax, isChild), 10, 64); err == nil && value.Int() > max {
				return errors.New(fmt.Sprint(name, " must not be greater than ", max))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if min, err := strconv.ParseUint(getTag(tag, tagMin, isChild), 10, 64); err == nil && value.Uint() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseUint(getTag(tag, tagMax, isChild), 10, 64); err == nil && value.Uint() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.Float32, reflect.Float64:
		if min, err := strconv.ParseFloat(getTag(tag, tagMin, isChild), 64); err == nil && value.Float() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseFloat(getTag(tag, tagMax, isChild), 64); err == nil && value.Float() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.String:
		if isEmpty, err := strconv.ParseBool(getTag(tag, tagIsEmpty, isChild)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(getTag(tag, tagMin, isChild)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " characters"))
		}
		if max, err := strconv.Atoi(getTag(tag, tagMax, isChild)); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " characters"))
		}
	case reflect.Map:
		if isEmpty, err := strconv.ParseBool(getTag(tag, tagIsEmpty, isChild)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(getTag(tag, tagMin, isChild)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " elements"))
		}
		if max, err := strconv.Atoi(getTag(tag, tagMax, isChild)); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " elements"))
		}
	case reflect.Slice:
		if isEmpty, err := strconv.ParseBool(getTag(tag, tagIsEmpty, isChild)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(getTag(tag, tagMin, isChild)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " elements"))
		}
		if max, err := strconv.Atoi(getTag(tag, tagMax, isChild)); err == nil && value.Len() > max {
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
		if isNil, err := strconv.ParseBool(getTag(tag, tagIsNil, isChild)); err == nil {
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

func getTag(tag reflect.StructTag, tagName string, child bool) string {
	var tagMap = map[string]string{
		tagMin:     tagChildMin,
		tagMax:     tagChildMax,
		tagIsEmpty: tagChildIsEmpty,
		tagIsNil:   tagChildIsNil,
	}

	if child {
		tagName = tagMap[tagName]
	}

	return tag.Get(tagName)
}
