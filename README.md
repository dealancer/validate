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
* `format` validator checks if a string in one of the following formats: `alpha`, `alnum`, `alpha_unicode`, `alnum_unicode`, `numeric`, `number`, `hexadecimal`, `hexcolor`, `rgb`, `rgba`, `hsl`, `hsla`, `email`, `url`, `uri`, `urn_rfc2141`, `file`, `base64`, `base64url`, `isbn`, `isbn10`, `isbn13`, `eth_addr`, `btc_addr`, `btc_addr_bech32`, `uuid`, `uuid3`, `uuid4`, `uuid5`, `ascii`, `ascii_print`, `datauri`, `latitude`, `longitude`, `ssn`, `ipv4`, `ipv6`, `ip`, `cidrv4`, `cidrv6`, `cidr`, `mac`, `hostname`, `hostname_rfc1123`, `fqdn`, `url_encoded`, `dir`.
* `[]` (brackets) are used to validate map keys.
* `>` (arrow) is used to validate values of maps, slices, arrays or to dereference a pointer.
* `;` (semicolon) is used to perform multiple validators using AND logic.
* `,` (commna) is used to specify multiple tokens for `one_of` validator.


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
    // Name should not be empty
    Name string `validate:"empty=false"`

    // Hosts should not be empty, Hosts values should be in the right format
    Hosts []string `validate:"empty=false > format=hostname"`

    // Username should be one of "joe", "ivan", or "li"
    Username string `validate:"one_of=joe,ivan,li"`

    // Password should be more than or equal to twelve characters
    Password *string `validate:"> min=12"`

    // Ssl (pointer) should not be nil
    Ssl *bool `validate:"nil=false"`

    // SslVerify (pointer) should not be nil
    SslVerify *bool `validate:"nil=false"`

    // Version should be between 5 and 8
    Version int `validate:"min=5; max=8"`
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