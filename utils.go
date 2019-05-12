package validate

import (
	"reflect"
	"strings"
)

func parseValidators(tag reflect.StructTag) map[string]string {
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

func getValidatorType(valType string, child bool) string {
	var valChildMap = map[string]string{
		valTypeChildMin:   valTypeMin,
		valTypeChildMax:   valTypeMax,
		valTypeChildEmpty: valTypeEmpty,
		valTypeChildNil:   valTypeNil,
		valTypeChildOneOf: valTypeOneOf,
	}

	if child {
		if valType, ok := valChildMap[valType]; ok {
			return valType
		}
		return ""
	}

	return valType
}

func parseTokens(str string) []interface{} {
	if strings.TrimSpace(str) == "" {
		return nil
	}

	tokenStrings := strings.Split(str, "|")
	tokens := make([]interface{}, len(tokenStrings))
	for i := range tokenStrings {
		tokens[i] = strings.TrimSpace(tokenStrings[i])
	}

	return tokens
}

func isOneOf(token interface{}, tokens []interface{}) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}

	return false
}
