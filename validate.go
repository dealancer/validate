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
)

var (
	errInvalidType = errors.New("not a struct pointer")
)

// Validate validates members in a struct
func Validate(ptr interface{}) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return errInvalidType
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return errInvalidType
	}

	for i := 0; i < t.NumField(); i++ {
		if err := validateField(t.Field(i), v.Field(i)); err != nil {
			return err
		}

	}

	return nil
}

func validateField(f reflect.StructField, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		if notEmpty, err := strconv.ParseBool(f.Tag.Get(tagNotEmpty)); err == nil && notEmpty && v.String() == "" {
			return errors.New(fmt.Sprint(f.Name, " must not be empty"))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if min, err := strconv.ParseInt(f.Tag.Get(tagMin), 10, 64); err == nil && v.Int() < min {
			return errors.New(fmt.Sprint(f.Name, " must not be less than ", min))
		}
		if max, err := strconv.ParseInt(f.Tag.Get(tagMax), 10, 64); err == nil && v.Int() > max {
			return errors.New(fmt.Sprint(f.Name, " must not be greater than ", max))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if min, err := strconv.ParseUint(f.Tag.Get(tagMin), 10, 64); err == nil && v.Uint() < min {
			return errors.New(fmt.Sprint(f.Name, " must not be less than ", min))
		}
		if max, err := strconv.ParseUint(f.Tag.Get(tagMax), 10, 64); err == nil && v.Uint() > max {
			return errors.New(fmt.Sprint(f.Name, " must not be greater than ", max))
		}
	}

	return nil
}
