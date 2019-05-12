# validate
[![Build Status](https://travis-ci.org/dealancer/validate.svg?branch=master)](https://travis-ci.org/dealancer/validate)
[![codecov](https://codecov.io/gh/dealancer/validate/branch/master/graph/badge.svg)](https://codecov.io/gh/dealancer/validate)
[![Go Report Card](https://goreportcard.com/badge/github.com/dealancer/validate)](https://goreportcard.com/report/github.com/dealancer/validate)
[![GoDoc](https://godoc.org/github.com/dealancer/validate?status.svg)](https://godoc.org/github.com/dealancer/validate)
[![GitHub release](https://img.shields.io/github/release/dealancer/validate.svg)](https://github.com/dealancer/validate/releases)
[![License](https://img.shields.io/github/license/dealancer/validate.svg)](./LICENSE)

**Validate** validates members of a Go struct.

## Types

This package supports a wide variety of types:

* Scalar types:
  * `int/8/16/32/64`, `uint/8/16/32/64`, `float32/64`
  * `uintptr`
  * `string`
* Complex types:
  * `map`
  * `slice`
* Aliased types:
  * `time.Duration`
  * e.g. `rune`, `type Enum string`
* Pointer types:
  * e.g, `*string`, `*int`
  
## Validators

This package supports following validators:

* `min`, `max`: works with scalar types (numbers, strings), aliased types, maps, and slices
* `one_of`: works with with scalar types (numbers, strings) and aliased types
* `empty`: works with strings, maps, and slices
* `nil`: works with pointers
* `child_min`, `child_max`, `child_one_of`, `child_empty`, `child_nil`: works with child elements of slices and referenced elements of pointers

## Installation

```
go get github.com/dealancer/validate
```

## Usage

```go
type Connection struct {
	Name      string   `validate:"empty=false"`
	Hosts     []string `validate:"empty=false,child_empty=false"`
	Username  string   `validate:"is_one_of=joe|ivan|li"`
	Password  *string  `validate:"child_min=12"`
	Ssl       *bool    `validate:"nil=false"`
	SslVerify *bool    `validate:"nil=false"`
	Version   int      `validate:"min=5,max=8"`

	XXX map[string]interface{} `validate:"empty=true"`
}
```

```go
connection := Connection{
	Username: "admin",
}

if err := validate.Validate(&connection); err != nil {
	panic(err)
}
```

## Unmarshalling YAML/JSON

This package can be used together with [github.com/creasty/defaults](http://github.com/creasty/defaults) for validating and providing default values for complex structs coming from YAML and JSON. This can be conveniently done by implementing `UnmarshalYAML` or `UnmarshalJSON` interfaces.

```go
func (this *Connection) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(this); err != nil {
		return err
	}

	type plain Connection
	if err := unmarshal((*plain)(this)); err != nil {
		return err
	}

	if err := validate.Validate(this); err != nil {
		return err
	}

	return nil
}
```
