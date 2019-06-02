# validate
[![Build Status](https://travis-ci.org/dealancer/validate.svg?branch=v1)](https://travis-ci.org/dealancer/validate)
[![codecov](https://codecov.io/gh/dealancer/validate/branch/v1/graph/badge.svg)](https://codecov.io/gh/dealancer/validate)
[![Go Report Card](https://goreportcard.com/badge/github.com/dealancer/validate)](https://goreportcard.com/report/github.com/dealancer/validate)
[![GoDoc](https://godoc.org/github.com/dealancer/validate?status.svg)](https://godoc.org/github.com/dealancer/validate)
[![GitHub release](https://img.shields.io/github/release/dealancer/validate.svg)](https://github.com/dealancer/validate/releases)
[![License](https://img.shields.io/github/license/dealancer/validate.svg)](./LICENSE)

Package **validate** validates fields of the Go struct recursively based on tags.
It provides powerful syntax to perform validation for substructs, maps, slices, arrays, and pointers.

Use this package to make sure that the content of the struct is in the format you need.
For example, **validate** package is useful when unmarshalling YAML or JSON.

## Types

This package supports a wide variety of types.

* Most of the built-in types:
  * `int`, `int8` `int16`, `int32`, `int64`
  * `uint`, `uint8`, `uint16`, `uint32`, `uint64`
  * `float32`, `float64`
  * `uintptr`
  * `string`
* Aliased types:
  * `time.Duration`
  * `byte` (`uint8`)
  * `rune` (`int32`)
  * e.g. `type Enum string`
* Complex types:
  * Struct
  * Map
  * Slice
  * Array
  * Pointer

## Validators

This package provides the following validators.

* `min` and `max` validators compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `empty` validator checks if a string, a map, a slice, or an array is (not) empty.
* `nil` validator checks if a pointer is (not) nil.
* `one_of` validator checks if a number or a string contains any of the given elements.
* `[]` (brackets) are used to validate map keys.
* `>` (arrow) is used to validate values of maps, slices, arrays or to dereference a pointer.

## Installation

1. Import validate package into your project.
   ```go
   import "github.com/dealancer/validate"
   ```
2. Additionally use `go get` when not using Go modules.
   ```
    go get github.com/dealancer/validate
    ```

## Usage

```go
type Connection struct {
	Name      string   `validate:"empty=false"`               // Name should not be empty
	Hosts     []string `validate:"empty=false > empty=false"` // Hosts should not be empty, Hosts values should not be empty
	Username  string   `validate:"one_of=joe,ivan,li"`        // Username should be one of "joe", "ivan", or "li"
	Password  *string  `validate:"> min=12"`                  // Password should be more than twelve characters
	Ssl       *bool    `validate:"nil=false"`                 // Ssl (pointer) should not be nil
	SslVerify *bool    `validate:"nil=false"`                 // SslVerify (pointer) should not be nil
	Version   int      `validate:"min=5; max=8"`              // Version should be between 5 and 8
}

type Connections struct {
	Connections []Connection `validate:"min=2"` // There should be at least two connections
}
```

```go
connections := Connections{
	Connections: []Connection{
		Connection{
			Username: "admin",
		},
	},
}

if err := validate.Validate(&connections); err != nil {
	panic(err)
}
```

See [GoDoc](https://godoc.org/github.com/dealancer/validate) for the complete reference.