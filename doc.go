/*
Package validate validates Go structs and types recursively based on tags.
It provides powerful syntax to perform validation for substructs, maps, slices, arrays, and pointers.
Package also allows to run custom validation methods.

Use this package to make sure that the content of the struct is in the format you need.
For example, **validate** package is useful when unmarshalling YAML or JSON.

This package supports most of the built-in types: int8, uint8, int16, uint16, int32,
uint32, int64, uint64, int, uint, uintptr, float32, float64 and aliased types:
time.Duration, byte (uint8), rune (int32).

Following validators are available: gt, lt, gte, lte, empty, nil, one_of, format.

Basic usage

Use validate tag to specify validators for fields of a struct.
If any of validators fail, validate.Validate returns an error.

	type S struct {
		i int    `validate:"gte=0"`        // Should be greater than or equal to 0
		s string `validate:"format=email"` // Should be in the email format
		b *bool  `validate:"nil=false"`    // Should not be nil
	}

	err := validate.Validate(S{
		i: -1,
		s: "",
		b: nil,
	})

	// err contains an error because n is less than 0, s is empty, and b is nil

Multiple validators

It is possible to specify multiple validators using & (ampersand) or | (vertical bar) operator.
& operator is used for logical AND, while | is used for logical OR.
& operator has a priority over | operator.

	type S struct {
	    // Check that the value is in the range of -20...-10 or 10...20
		field int `validate:"gte=-20 & lte=-10 | gte=10 & lte=20"`
	}

Slice and array validation

You can use a regular syntax to validate a slice/array. To validate slice/array values, specify validators after an arrow character.

	type S struct {
		// Check that the slice is not empty and slice values are either 1 or -1
		field []int `validate:"empty=false > one_of=1,-1"`
	}

Map validation

You can use a regular syntax to validate a map. To validate map keys, specify validators inside brackets.
To validate map values, specify validators after an arrow character.

	type S struct {
		// Check that the map contains at least two elements, map keys are not empty, and map values are between 0 and 10
		field map[string]int `validate:"gte=2 [empty=false] > gte=0 & lte=10"`
	}

Pointer validation

You can use a regular syntax to validate a pointer. To dereference a pointer, specify validators after an arrow character.

	type S struct {
		// Check that the pointer is not nil and the number is between 0 and 10
		field *int `validate:"nil=false > gte=0 & lte=10"`
	}

Nested struct validation

You can validate a nested struct with regular syntax.

	type A struct {
		// Check that the number is greater than or equal to 0
		a int `validate:"gte=0"`
	}

	type B struct {
		A
		// Check that the number is greater than or equal to 0
		b int `validate:"gte=0"`
	}

Substruct validation

You can validate a substruct with regular syntax.

	type A struct {
		// Check that the number is greater than or equal to 0
		field int `validate:"gte=0"`
	}

	type B struct {
		a A
		// Check that the number is greater than or equal to 0
		field int `validate:"gte=0"`
	}

Deep validation

You can use brackets and arrow syntax to deep to as many levels as you need.

	type A struct {
		field int `validate:"gte=0 & lte=10"`
	}

	type B struct {
		field []map[*string]*A `validate:"gte=1 & lte=2 | eq=4 > empty=false [nil=false > empty=false] > nil=false"`
	}

	// gte=1 & lte=2 | eq=4 will be applied to the array
	// empty=false will be applied to the map
	// nil=false > empty=false will be applied to the map key (pointer and string)
	// nil=false will be applied to the map value
	// gte=0 & lte=10 will be applied to the A.field

Custom validation

You can specify custom validation method.
Custom validation also works for a substuct, if a substruct is defined in an exported field.

	type S struct {
		field        int
	}

	func (s S) Validate() error {
		if s.field <= 0 {
			return errors.New("field should be positive")
		}
		return nil
	}

Handling errors

Validate method returns two types of errors: ErrorSyntax and ErrorValidation.
You can handle an error type using switch syntax.

	type S struct {
		field *string `validate:"empty=false"`
	}

	var err error

	if err = validate.Validate(S{nil}); err != nil {
		switch err.(type) {
		case validate.ErrorSyntax:
			// Handle syntax error
		case validate.ErrorValidation:
			// Handle validation error
		default:
			// Handle other errors
		}
	}
*/
package validate
