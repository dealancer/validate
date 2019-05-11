# validate
[![GitHub release](https://img.shields.io/github/release/dealancer/validate.svg)](https://github.com/dealancer/validate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/dealancer/validate)](https://goreportcard.com/report/github.com/dealancer/validate)
[![GoDoc](https://godoc.org/github.com/dealancer/validate?status.svg)](https://godoc.org/github.com/dealancer/validate)
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
  
## Validation tags

This package supports following tags:

* `is_empty`: works with strings, maps, and slices
* `is_nil`: works with pointers
* `min`, `max`: works with numbers, strings, maps, and slices
* `child_is_empty`, `child_is_nil`, `child_min`, `child_max`: works with child elementts of slices and referenced elements of pointers

## Installation

```
go get github.com/dealancer/validate
```

## Usage

```go
type Connection struct {
	Name      string   `is_empty:"false"`
	Hosts     []string `is_empty:"false" child_is_empty:"false"`
	Username  string   `is_empty:"false"`
	Password  *string  `child_min:"12"`
	Ssl       *bool    `is_nil:"false"`
	SslVerify *bool    `is_nil:"false""`
	Version   int      `min:"5" max:"8"`

	XXX map[string]interface{} `is_empty:"true"`
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

This package can be used together with [github.com/creasty/defaults](http://github.com/creasty/defaults) for validating and providing default values for complex structs coming from YAML and JSON. This can be conveniently by implementing `UnmarshalYAML` or `UnmarshalJSON` interfaces.

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
