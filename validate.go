// Package validate validates Go structs and types recursively based on tags.
// It provides powerful syntax to perform validation for substructs, maps, slices, arrays, and pointers.
// Package also allows to run custom validation methods.
//
// Use this package to make sure that the content of the struct is in the format you need.
// For example, **validate** package is useful when unmarshalling YAML or JSON.
//
// This package supports most of the built-in types: int8, uint8, int16, uint16, int32,
// uint32, int64, uint64, int, uint, uintptr, float32, float64 and aliased types:
// time.Duration, byte (uint8), rune (int32).
//
// Following validators are available: gt, lt, gte, lte, empty, nil, one_of, format.
//
// Basic usage
//
// Use validate tag to specify validators for fields of a struct.
// If any of validators fail, validate.Validate returns an error.
//
//  type S struct {
//  	i int    `validate:"gte=0"`        // Should be greater than or equal to 0
//  	s string `validate:"format=email"` // Should be in the email format
//  	b *bool  `validate:"nil=false"`    // Should not be nil
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
// It is possible to specify multiple validators using & (ampersand) or | (vertical bar) operator.
// & operator is used for logical AND, while | is used for logical OR.
// & operator has a priority over | operator.
//
//  type S struct {
//      // Check that the value is in the range of -20...-10 or 10...20
//  	field int `validate:"gte=-20 & lte=-10 | gte=10 & lte=20"`
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
//  	field map[string]int `validate:"gte=2 [empty=false] > gte=0 & lte=10"`
//  }
//
// Pointer validation
//
// You can use a regular syntax to validate a pointer. To dereference a pointer, specify validators after an arrow character.
//
//  type S struct {
//  	// Check that the pointer is not nil and the number is between 0 and 10
//  	field *int `validate:"nil=false > gte=0 & lte=10"`
//  }
//
// Nested struct validation
//
// You can validate a nested struct with regular syntax.
//
//  type A struct {
//  	// Check that the number is greater than or equal to 0
//  	a int `validate:"gte=0"`
//  }
//
//  type B struct {
//  	A
//  	// Check that the number is greater than or equal to 0
//  	b int `validate:"gte=0"`
//  }
//
// Substruct validation
//
// You can validate a substruct with regular syntax.
//
//  type A struct {
//  	// Check that the number is greater than or equal to 0
//  	field int `validate:"gte=0"`
//  }
//
//  type B struct {
//  	a A
//  	// Check that the number is greater than or equal to 0
//  	field int `validate:"gte=0"`
//  }
//
// Deep validation
//
// You can use brackets and arrow syntax to deep to as many levels as you need.
//
//  type A struct {
//  	field int `validate:"gte=0 & lte=10"`
//  }
//
//  type B struct {
//  	field []map[*string]*A `validate:"gte=1 & lte=2 | eq=4 > empty=false [nil=false > empty=false] > nil=false"`
//  }
//
//  // gte=1 & lte=2 | eq=4 will be applied to the array
//  // empty=false will be applied to the map
//  // nil=false > empty=false will be applied to the map key (pointer and string)
//  // nil=false will be applied to the map value
//  // gte=0 & lte=10 will be applied to the A.field
//
// Custom validation
//
// You can specify custom validation method.
// Custom validation also works for a substuct, if a substruct is defined in an exported field.
//
//  type S struct {
//  	field        int
//  }
//
//  // Make sure reciever is value, otherwise it won't work
//  func (s S) Validate() error {
//  	if s.field <= 0 {
//  		return errors.New("field should be positive")
//  	}
//
//  	return nil
//  }
//
package validate

import (
	"reflect"
	"regexp"
	"strings"
)

// MasterTag is the main validation tag.
const MasterTag = "validate"

// CustomValidator is an interface for a validated struct.
type CustomValidator interface {

	// Validate is a custom validation function.
	// Validate does not work when the reciever is a reference.
	// Validate does not work for nested types obtained from unexported field.
	Validate() error
}

// Validate validates fields of a struct.
// It accepts a struct or a struct pointer as a parameter.
// It returns an error if a struct does not validate or nil if there are no validation errors.
//
//  err := validate.Validate(struct {
//  	field time.Duration `validate:"gte=0s"`
//  }{
//  	field: -time.Second,
//  })
//
//  // err contains an error
func Validate(element interface{}) error {
	value := reflect.ValueOf(element)

	return validateField(value, "", "")
}

// validateField validates a struct field
func validateField(value reflect.Value, fieldName string, validators string) error {
	kind := value.Kind()

	// Get validator type Map
	validatorTypeMap := getValidatorTypeMap()

	// Get validators
	keyValidators, valueValidators, validators := splitValidators(validators)
	validatorsOr := parseValidators(valueValidators)

	// Call a custom validator
	if value.CanInterface() {
		if customValidator, ok := value.Interface().(CustomValidator); ok {
			err := customValidator.Validate()
			if err != nil {
				return err
			}
		}
	}

	// Perform validators
	var err error
	for _, validatorsAnd := range validatorsOr {
		for _, validator := range validatorsAnd {
			if validatorFunc, ok := validatorTypeMap[validator.Type]; ok {
				if err = validatorFunc(value, fieldName, validator.Value); err != nil {
					break
				}
			}
		}
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
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

// validateStruct validates a struct
func validateStruct(value reflect.Value) error {
	typ := value.Type()

	// Iterate over struct fields
	for i := 0; i < typ.NumField(); i++ {
		validators := getValidators(typ.Field(i).Tag)
		fieldName := typ.Field(i).Name
		if err := validateField(value.Field(i), fieldName, validators); err != nil {
			return err
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

// parseValidator2 parses validators into the slice of slices.
// First slice acts as AND logic, second array acts as OR logic.
func parseValidators(validators string) (validatorsOr [][]validator) {
	regexpType := regexp.MustCompile(`[[:alnum:]_]+`)
	regexpValue := regexp.MustCompile(`[^=\s]+[^=]*[^=\s]+|[^=\s]+`)

	entriesOr := strings.Split(validators, "|")
	validatorsOr = make([][]validator, 0, len(entriesOr))
	for _, entryOr := range entriesOr {
		entriesAnd := strings.Split(entryOr, "&")
		validatorsAnd := make([]validator, 0, len(entriesAnd))
		for _, entryOr := range entriesAnd {
			entries := strings.Split(entryOr, "=")
			if len(entries) > 0 {
				t := regexpType.FindString(entries[0])
				v := ""
				if len(entries) == 2 {
					v = regexpValue.FindString(entries[1])
				}
				if len(t) > 0 {
					validatorsAnd = append(validatorsAnd, validator{ValidatorType(t), v})
				}
			}
		}
		if len(validatorsAnd) > 0 {
			validatorsOr = append(validatorsOr, validatorsAnd)
		}
	}

	return
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
