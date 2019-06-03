// Package validate validates fields of the Go struct recursively based on tags.
// It provides powerful syntax to perform validation for substructs, maps, slices, arrays, and pointers.
//
// Use this package to make sure that the content of the struct is in the format you need.
// For example, **validate** package is useful when unmarshalling YAML or JSON.
//
// This package supports most of the built-in types: int8, uint8, int16, uint16, int32,
// uint32, int64, uint64, int, uint, uintptr, float32, float64 and aliased types:
// time.Duration, byte (uint8), rune (int32).
//
// Following validators are available: min, max, empty, nil, one_of.
//
// Basic usage
//
// Use validate tag to specify validators for fields of a struct.
// If any of validators fail, validate.Validate returns an error.
//
//  type S struct {
//  	i int    `validate:"min=0"`       // Should be greater than or equal to 0
//  	s string `validate:"empty=false"` // Should not be empty
//  	b *bool  `validate:"nil=false"`   // Should not be nil
//  }
//
//  err := validate.Validate(S{
//  	i: -1,
//  	s: "",
//  	b: nil,
//  })
//
//  // err contains an error because n is less than 0, s is empty, and b is nil
//
// Multiple validators
//
// It is possible to specify multiple validators using a semicolon character.
//
//  type S struct {
//  	field int `validate:"min=0; max=10"`
//  }
//
// Slice and array validation
//
// You can use a regular syntax to validate a slice/array. To validate slice/array values, specify validators after an arrow character.
//
//  type S struct {
//  	// Check that the slice is not empty and slice values are either 1 or -1
//  	field []int `validate:"empty=false > one_of=1,-1"`
//  }
//
// Map validation
//
// You can use a regular syntax to validate a map. To validate map keys, specify validators inside brackets.
// To validate map values, specify validators after an arrow character.
//
//  type S struct {
//  	// Check that the map contains at least two elements, map keys are not empty, and map values are between 0 and 10
//  	field map[string]int `validate:"min=2 [empty=false] > min=0; max=10"`
//  }
//
// Pointer validation
//
// You can use a regular syntax to validate a pointer. To dereference a pointer, specify validators after an arrow character.
//
//  type S struct {
//  	// Check that the pointer is not nil and the number is between 0 and 10
//  	field *int `validate:"nil=false > min=0; max=10"`
//  }
//
// Nested struct validation
//
// You can validate a nested struct with regular syntax.
//
//  type A struct {
//  	// Check that the number is greater than or equal to 0
//  	a int `validate:"min=0"`
//  }
//
//  type B struct {
//  	A
//  	// Check that the number is greater than or equal to 0
//  	b int `validate:"min=0"`
//  }
//
// Substruct validation
//
// You can validate a substruct with regular syntax.
//
//  type A struct {
//  	// Check that the number is greater than or equal to 0
//  	field int `validate:"min=0"`
//  }
//
//  type B struct {
//  	a A
//  	// Check that the number is greater than or equal to 0
//  	field int `validate:"min=0"`
//  }
//
// Deep validation
//
// You can use brackets and arrow syntax to deep to as many levels as you need.
//
//  type A struct {
//  	field int `validate:"min=0; max=10"`
//  }
//
//  type B struct {
//  	field []map[*string]*A `validate:"min=1; max=2 > empty=false [nil=false > empty=false] > nil=false"`
//  }
//
//  // min=1, max=2 will be applied to the array
//  // empty=false will be applied to the map
//  // nil=false > empty=false will be applied to the map key (pointer and string)
//  // nil=false will be applied to the map value
//  // min=0, max=10 will be applied to the A.field
//
package validate

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// MasterTag is the main validation tag.
const MasterTag = "validate"

// Validate validates fields of a struct.
// It accepts a struct or a struct pointer as a parameter.
// It returns an error if a struct does not validate or nil if there are no validation errors.
//
//  err := validate.Validate(struct {
//  	field time.Duration `validate:"min=0s"`
//  }{
//  	field: -time.Second,
//  })
//
//  // err contains an error
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
		validators := getValidators(typ.Field(i).Tag)
		fieldName := typ.Field(i).Name
		if err := validateField(value.Field(i), fieldName, validators); err != nil {
			return err
		}
	}

	return nil
}

// validateField validates a struct field
func validateField(value reflect.Value, fieldName string, validators string) error {
	kind := value.Kind()

	// Get validators
	keyValidators, valueValidators, validators := splitValidators(validators)
	valueValidatorMap := parseValidators(valueValidators)

	// Perform validators
	validatorTypeMap := getValidatorTypeMap()
	for validatorType, validator := range valueValidatorMap {
		if validatorFunc, ok := validatorTypeMap[validatorType]; ok {
			if err := validatorFunc(value, fieldName, validator); err != nil {
				return err
			}
		}
	}

	// Dive one level deep into arrays and pointers
	switch kind {
	case reflect.Struct:
		if err := validateStruct(value); err != nil {
			return err
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			if err := validateField(key, fieldName, keyValidators); err != nil {
				return err
			}
			if err := validateField(value.MapIndex(key), fieldName, validators); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if err := validateField(value.Index(i), fieldName, validators); err != nil {
				return err
			}
		}
	case reflect.Ptr:
		if !value.IsNil() {
			if err := validateField(value.Elem(), fieldName, validators); err != nil {
				return err
			}
		}
	}

	return nil
}

// getValidators gets validators
func getValidators(tag reflect.StructTag) string {
	return tag.Get(MasterTag)
}

// splitValidators splits validators into key validators, value validators and remaning validators of the next level
func splitValidators(validators string) (keyValidators string, valValidators string, remaningValidators string) {
	bracket := 0
	bracketStart := 0
	bracketEnd := -1

	i := 0
loop:
	for ; i < len(validators); i++ {
		switch validators[i] {
		case '>':
			if bracket == 0 {
				break loop
			}
		case '[':
			if bracket == 0 {
				bracketStart = i
			}
			bracket++
		case ']':
			bracket--
			if bracket == 0 {
				bracketEnd = i
			}
		}
	}

	if bracketStart <= len(validators) {
		valValidators += validators[:bracketStart]
	}
	if bracketEnd+1 <= len(validators) {
		if valValidators != "" {
			valValidators += " "
		}
		valValidators += validators[bracketEnd+1 : i]
	}
	if bracketStart+1 <= len(validators) && bracketEnd >= 0 && bracketStart+1 <= bracketEnd {
		keyValidators = validators[bracketStart+1 : bracketEnd]
	}
	if i+1 <= len(validators) {
		remaningValidators = validators[i+1:]
	}

	return
}

// parseValidators parses validators into the hash map
func parseValidators(validators string) (validatorMap map[ValidatorType]string) {
	validatorMap = make(map[ValidatorType]string)

	r := regexp.MustCompile(`([[:alnum:]_\s]+)=?([^=;]*);?`)

	entries := r.FindAllStringSubmatch(validators, -1)

	for _, e := range entries {
		n := strings.TrimSpace(e[1])
		v := strings.TrimSpace(e[2])

		if n != "" {
			validatorMap[ValidatorType(n)] = v
		}
	}

	return validatorMap
}

// parseTokens parses tokens into array
func parseTokens(str string) []interface{} {
	tokenStrings := strings.Split(str, ",")
	tokens := make([]interface{}, 0, len(tokenStrings))

	for i := range tokenStrings {
		token := strings.TrimSpace(tokenStrings[i])
		if token != "" {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// tokenOneOf check if a token is one of tokens
func tokenOneOf(token interface{}, tokens []interface{}) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}

	return false
}
