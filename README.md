# validate
Validate members in a struct.

## Types

This package supports for a wide variety of types:

* Scalar types
  * `int/8/16/32/64`, `uint/8/16/32/64`, `float32/64`
  * `uintptr`
  * `string`
* Complex types
  * `map`
  * `slice`
* Aliased types
  * time.Duration
  * e.g. type Enum string
* Pointer types
  * e.g, `*string`, `*int`
  
## Validation tags

This package supports following tags:

* `is_empty`: works with strings, maps, and slices
* `is_nil`: works with pointers
* `min`: works with numbers, strings, maps, slices
* `max`: works with numbers, strings, maps, slices

## Installation

```
go get github.com/dealancer/validate
```

## Usage

```
type Connection struct {
	Name      string   `is_empty:"false"`
	Hosts     []string `is_empty:"false"`
	Username  string   `is_empty:"false"`
	Password  string   `min:"12"`
	Version   int      `min:"5" max:"8"`
	Ssl       *bool    `is_nil:"false"`
	SslVerify *bool    `is_nil:"false"`

	XXX map[string]interface{} `is_empty:"true"`
}
```

```
connection := Connection{
	Username: "admin",
	Password: "password",
}

if err := validate.Validate(&connection); err != nil {
	panic(err)
}
```

## Unmarshalling YAML/JSON

This package can be used togeather with [github.com/creasty/defaults](http://github.com/creasty/defaults) for validating and providing default values for complex structs coming from YAML and JSON. This can be conviniently by implementing `UnmarshalYAML` or `UnmarshalJSON` interfaces.

```
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
