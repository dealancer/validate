package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagMin     = "min"
	tagMax     = "max"
	tagIsEmpty = "is_empty"
	tagIsNil   = "is_nil"
)

var (
	errInvalidType = errors.New("not a struct pointer")
)

// Validate validates members of a struct
func Validate(ptr interface{}) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return errInvalidType
	}

	value := reflect.ValueOf(ptr).Elem()
	typ := value.Type()

	if typ.Kind() != reflect.Struct {
		return errInvalidType
	}

	for i := 0; i < typ.NumField(); i++ {
		if err := validateField(value.Field(i), value.Field(i).Kind(), typ.Field(i).Name, typ.Field(i).Tag); err != nil {
			return err
		}
	}

	return nil
}

func validateField(value reflect.Value, kind reflect.Kind, name string, tag reflect.StructTag) error {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if min, err := strconv.ParseInt(tag.Get(tagMin), 10, 64); err == nil && value.Int() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseInt(tag.Get(tagMax), 10, 64); err == nil && value.Int() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if min, err := strconv.ParseUint(tag.Get(tagMin), 10, 64); err == nil && value.Uint() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseUint(tag.Get(tagMax), 10, 64); err == nil && value.Uint() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.Float32, reflect.Float64:
		if min, err := strconv.ParseFloat(tag.Get(tagMin), 64); err == nil && value.Float() < min {
			return errors.New(fmt.Sprint(name, " must not be less than ", min))
		}
		if max, err := strconv.ParseFloat(tag.Get(tagMax), 64); err == nil && value.Float() > max {
			return errors.New(fmt.Sprint(name, " must not be greater than ", max))
		}
	case reflect.String, reflect.Map, reflect.Slice:
		if isEmpty, err := strconv.ParseBool(tag.Get(tagIsEmpty)); err == nil {
			if isEmpty && value.Len() > 0 {
				return errors.New(fmt.Sprint(name, " must be empty"))
			} else if !isEmpty && value.Len() == 0 {
				return errors.New(fmt.Sprint(name, " must not be empty"))
			}
		}
		if min, err := strconv.Atoi(tag.Get(tagMin)); err == nil && value.Len() < min {
			return errors.New(fmt.Sprint(name, " must not contain less than ", min, " elements"))
		}
		if max, err := strconv.Atoi(tag.Get(tagMax)); err == nil && value.Len() > max {
			return errors.New(fmt.Sprint(name, " must not contain more than ", max, " elements"))
		}
	case reflect.Ptr:
		if isNil, err := strconv.ParseBool(tag.Get(tagIsNil)); err == nil {
			if isNil && !value.IsNil() {
				return errors.New(fmt.Sprint(name, " must be nil"))
			} else if !isNil && value.IsNil() {
				return errors.New(fmt.Sprint(name, " must not be nil"))
			}
		}
		if !value.IsNil() {
			return validateField(value.Elem(), value.Elem().Kind(), name, tag)
		}
	}

	return nil
}
