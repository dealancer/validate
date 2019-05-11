package validate

import (
	"testing"
	"time"
)

func TestMinTagForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `min:"0s"`
	}{
		field: -time.Second,
	}) {
		t.Errorf("min tag does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `min:"-1s"`
	}{
		field: -time.Minute,
	}) {
		t.Errorf("min tag does not validate for time.Duratuon")
	}
}

func TestMaxTagForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `max:"0s"`
	}{
		field: time.Second,
	}) {
		t.Errorf("max tag does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `max:"1s"`
	}{
		field: time.Minute,
	}) {
		t.Errorf("max tag does not validate for time.Duratuon")
	}
}

func TestMinTagForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `min:"0"`
	}{
		field: -1,
	}) {
		t.Errorf("min tag does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `min:"0"`
	}{
		field: -1,
	}) {
		t.Errorf("min tag does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `min:"0"`
	}{
		field: -1,
	}) {
		t.Errorf("min tag does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `min:"0"`
	}{
		field: -1,
	}) {
		t.Errorf("min tag does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `min:"0"`
	}{
		field: -1,
	}) {
		t.Errorf("min tag does not validate for int64")
	}
}

func TestMaxTagForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for int64")
	}
}
func TestMinTagForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `min:"10"`
	}{
		field: 9,
	}) {
		t.Errorf("min tag does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `min:"10"`
	}{
		field: 9,
	}) {
		t.Errorf("min tag does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `min:"10"`
	}{
		field: 9,
	}) {
		t.Errorf("min tag does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `min:"10"`
	}{
		field: 9,
	}) {
		t.Errorf("min tag does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `min:"10"`
	}{
		field: 9,
	}) {
		t.Errorf("min tag does not validate for uint64")
	}
}

func TestMaxTagForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `max:"0"`
	}{
		field: 1,
	}) {
		t.Errorf("max tag does not validate for uint64")
	}
}

func TestMinTagForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `min:"0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("min tag does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `min:"0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("min tag does not validate for flaot64")
	}
}

func TestMaxTagForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `max:"0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("max tag does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `max:"0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("max tag does not validate for flaot64")
	}
}
