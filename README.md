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

* `eq` (equals) validator compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `ne` (not equals) validators compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `gt` (greater than) validator compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `lt` (less than) validators compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `gte` (greater than or equal to) validator compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `lte` (less than or equal to) validators compare a numeric value of a number or compare a count of elements in a string, a map, a slice, or an array.
* `empty` validator checks if a string, a map, a slice, or an array is (not) empty.
* `nil` validator checks if a pointer is (not) nil.
* `one_of` validator checks if a number or a string contains any of the given elements.
* `format` validator checks if a string in one of the following formats: `alpha`, `alnum`, `alpha_unicode`, `alnum_unicode`, `numeric`, `number`, `hexadecimal`, `hexcolor`, `rgb`, `rgba`, `hsl`, `hsla`, `email`, `url`, `uri`, `urn_rfc2141`, `file`, `base64`, `base64url`, `isbn`, `isbn10`, `isbn13`, `eth_addr`, `btc_addr`, `btc_addr_bech32`, `uuid`, `uuid3`, `uuid4`, `uuid5`, `ascii`, `ascii_print`, `datauri`, `latitude`, `longitude`, `ssn`, `ipv4`, `ipv6`, `ip`, `cidrv4`, `cidrv6`, `cidr`, `mac`, `hostname`, `hostname_rfc1123`, `fqdn`, `url_encoded`, `dir`.

## Operators

Following operators are used. There are listed in the descending order of their precedence.

* `[]` (brackets) are used to validate map keys.
* `>` (greater-than sign) is used to validate values of maps, slices, arrays or to dereference a pointer.
* `&` (ampersand) is used to perform multiple validators using AND logic.
* `|` (vertical bar) is used to perform multiple validators using OR logic.
* `=` (equal sign) is used to separate validator type from value.
* `,` (comma) is used to specify multiple tokens for a validator (e.g. `one_of`).

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
    Password *string `validate:"> gte=12"`

    // Ssl (pointer) should not be nil
    Ssl *bool `validate:"nil=false"`

    // SslVerify (pointer) should not be nil
    SslVerify *bool `validate:"nil=false"`

    // Version should be between 5 and 8, or 9
    Version int `validate:"gte=5 & lte=8 | eq=9"`
}

type Connections struct {
	Connections []Connection `validate:"gte=2"` // There should be at least two connections
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