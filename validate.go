package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagMin      = "validate_min"
	tagMax      = "validate_max"
	tagNotEmpty = "validate_not_empty"
	tagNotNil   = "validate_not_nil"
)

var (
	errInvalidType = errors.New("not a struct pointer")
)

// Validate validates members in a struct
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
	case reflect.String:
		if notEmpty, err := strconv.ParseBool(tag.Get(tagNotEmpty)); err == nil && notEmpty && value.String() == "" {
			return errors.New(fmt.Sprint(name, " must not be empty"))
		}
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
	case reflect.Slice:
		if notEmpty, err := strconv.ParseBool(tag.Get(tagNotEmpty)); err == nil && notEmpty && value.Len() == 0 {
			return errors.New(fmt.Sprint(name, " must not be empty"))
		}
		for i := 0; i < value.Len(); i++ {
			if err := validateField(value.Index(i), value.Index(i).Kind(), name, tag); err != nil {
				return err
			}
		}
	case reflect.Ptr:
		if value.IsNil() {
			if notNil, err := strconv.ParseBool(tag.Get(tagNotNil)); err == nil && notNil {
				return errors.New(fmt.Sprint(name, " must not be nil"))
			}
		}
		return validateField(value.Elem(), value.Elem().Kind(), name, tag)
	}

	return nil
}
