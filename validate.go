package validate

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

const (
	masterTag = "validate"

	valTypeMin   = "min"
	valTypeMax   = "max"
	valTypeEmpty = "empty"
	valTypeNil   = "nil"
	valTypeOneOf = "one_of"
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

	// Perform validators
	_, valValidators, validators := splitValidators(validators)
	valValidatorsMap := parseValidators(valValidators)

	for valType, validator := range valValidatorsMap {
		var err error

		switch valType {
		case valTypeMin:
			err = validateMin(value, fieldName, validator)
		case valTypeMax:
			err = validateMax(value, fieldName, validator)
		case valTypeEmpty:
			err = validateEmpty(value, fieldName, validator)
		case valTypeNil:
			err = validateNil(value, fieldName, validator)
		case valTypeOneOf:
			err = validateOneOf(value, fieldName, validator)
		}

		if err != nil {
			return err
		}
	}

	// Dive one level deep into arrays and pointers
	switch kind {
	case reflect.Slice:
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
	return tag.Get(masterTag)
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
func parseValidators(validators string) (valMap map[string]string) {
	valMap = make(map[string]string)

	r, err := regexp.Compile(`([[:alnum:]_\s]+)=?([^=;]*);?`)
	if err != nil {
		return
	}

	entries := r.FindAllStringSubmatch(validators, -1)

	for _, e := range entries {
		n := strings.TrimSpace(e[1])
		v := strings.TrimSpace(e[2])

		if n != "" {
			valMap[n] = v
		}
	}

	return valMap
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

// isOneOf check if a token is one of tokens
func isOneOf(token interface{}, tokens []interface{}) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}

	return false
}
