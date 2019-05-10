package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Validate validates members in a struct
func Validate(ptr interface{}) error {
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		return nil
	}

	v := reflect.ValueOf(ptr).Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			if notEmpty, err := strconv.ParseBool(t.Field(i).Tag.Get("not_empty")); err == nil && notEmpty && v.Field(i).String() == "" {
				return errors.New(fmt.Sprint(t.Field(i).Name, " must not be empty"))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if min, err := strconv.ParseInt(t.Field(i).Tag.Get("min"), 10, 64); err == nil && v.Field(i).Int() < min {
				return errors.New(fmt.Sprint(t.Field(i).Name, " must not be less than ", min))
			}
			if max, err := strconv.ParseInt(t.Field(i).Tag.Get("max"), 10, 64); err == nil && v.Field(i).Int() > max {
				return errors.New(fmt.Sprint(t.Field(i).Name, " must not be greater than ", max))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if min, err := strconv.ParseUint(t.Field(i).Tag.Get("min"), 10, 64); err == nil && v.Field(i).Uint() < min {
				return errors.New(fmt.Sprint(t.Field(i).Name, " must not be less than ", min))
			}
			if max, err := strconv.ParseUint(t.Field(i).Tag.Get("max"), 10, 64); err == nil && v.Field(i).Uint() > max {
				return errors.New(fmt.Sprint(t.Field(i).Name, " must not be greater than ", max))
			}
		}
	}

	return nil
}
