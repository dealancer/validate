package validate

import (
	"reflect"
	"testing"
	"time"
)

func TestSplitValidators(t *testing.T) {
	key, val, validators := "", "", ""

	key, val, validators = splitValidators("")
	if key != "" || val != "" || validators != "" {
		t.Errorf("splitValidators incorrectly splits validators")
	}

	key, val, validators = splitValidators("[]>")
	if key != "" || val != "" || validators != "" {
		t.Errorf("splitValidators incorrectly splits validators")
	}

	key, val, validators = splitValidators(" [ ] > ")
	if key != " " || val != "   " || validators != " " {
		t.Errorf("splitValidators incorrectly splits validators")
	}

	key, val, validators = splitValidators("[[>]]>>")
	if key != "[>]" || val != "" || validators != ">" {
		t.Errorf("splitValidators incorrectly splits validators")
	}

	key, val, validators = splitValidators("val_a=a val_b=b")
	if key != "" || val != "val_a=a val_b=b" || validators != "" {
		t.Errorf("splitValidators incorrectly splits validators")
	}

	key, val, validators = splitValidators("val_a=a val_b=b [[val_c=c] > val_d=d] val_e=e > val_f=f [val_g=g] > val_h=h")
	if key != "[val_c=c] > val_d=d" || val != "val_a=a val_b=b   val_e=e " || validators != " val_f=f [val_g=g] > val_h=h" {
		t.Errorf("splitValidators incorrectly splits validators")
	}
}

func TestParseValidators(t *testing.T) {
	var valMap map[string]string

	valMap = parseValidators(";,;,;")
	if !reflect.DeepEqual(valMap, map[string]string{}) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	valMap = parseValidators("")
	if !reflect.DeepEqual(valMap, map[string]string{}) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	valMap = parseValidators("val_a=a")
	if !reflect.DeepEqual(valMap, map[string]string{
		"val_a": "a",
	}) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	valMap = parseValidators(" val  ;val_a=a;val_1 = 1  ;  val_b = b , c_d_ , 1.0 ")
	if !reflect.DeepEqual(valMap, map[string]string{
		"val":   "",
		"val_a": "a",
		"val_1": "1",
		"val_b": "b , c_d_ , 1.0",
	}) {
		t.Errorf("parseValidators incorrectly parses validators")
	}
}

func TestParseTokens(t *testing.T) {
	var tokens, res []interface{}

	tokens = parseTokens("")
	if !reflect.DeepEqual(tokens, make([]interface{}, 0)) {
		t.Errorf("parseTokens incorrectly parses validators")
	}

	tokens = parseTokens(" ,,  , ")
	if !reflect.DeepEqual(tokens, make([]interface{}, 0)) {
		t.Errorf("parseTokens incorrectly parses validators")
	}

	tokens = parseTokens("a, b, c")
	res = make([]interface{}, 3)
	res[0] = "a"
	res[1] = "b"
	res[2] = "c"
	if !reflect.DeepEqual(tokens, res) {
		t.Errorf("parseTokens incorrectly parses validators")
	}
}

func TestBasic(t *testing.T) {
	v := 1
	if nil == Validate(v) {
		t.Errorf("validate validates int type")
	}
	if nil == Validate(&v) {
		t.Errorf("validate validates &int type")
	}

	s := ""
	if nil == Validate(s) {
		t.Errorf("validate validates string type")
	}
	if nil == Validate(&s) {
		t.Errorf("validate validates &string type")
	}

	m := map[string]string{
		"a": "a",
	}
	if nil == Validate(m) {
		t.Errorf("validate validates map type")
	}
	if nil == Validate(m) {
		t.Errorf("validate validates &map type")
	}

	sl := []string{
		"a", "b",
	}
	if nil == Validate(sl) {
		t.Errorf("validate validates slice type")
	}
	if nil == Validate(&sl) {
		t.Errorf("validate validates slice type")
	}

	st := struct {
		field int
	}{
		field: 1,
	}
	if nil != Validate(st) {
		t.Errorf("validate does not validate struct type")
	}
	if nil != Validate(&st) {
		t.Errorf("validate does not validate struct pointer type")
	}

	stFail := struct {
		field int `validate:"max=0"`
	}{
		field: 1,
	}
	if nil == Validate(stFail) {
		t.Errorf("validate does not validate struct type")
	}
	if nil == Validate(&stFail) {
		t.Errorf("validate does not validate struct pointer type")
	}

	stAnotherFail := struct {
		a     int
		b     int
		field int `validate:"max=0"`
		c     int
		d     int
	}{
		field: 1,
	}
	if nil == Validate(stAnotherFail) {
		t.Errorf("validate does not validate struct type")
	}
	if nil == Validate(&stAnotherFail) {
		t.Errorf("validate does not validate struct pointer type")
	}
}

func TestMultiVal(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"min=0;max=10"`
	}{
		field: -1,
	}) {
		t.Errorf("multiple validators does not validate")
	}

	if nil == Validate(struct {
		field int `validate:"min=0;max=10"`
	}{
		field: 11,
	}) {
		t.Errorf("multiple validators does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"min=0;max=10"`
	}{
		field: 5,
	}) {
		t.Errorf("multiple validators does not validate")
	}

	if nil == Validate(struct {
		field int `validate:"min=1;max=-1"`
	}{
		field: 0,
	}) {
		t.Errorf("multiple validators does not validate")
	}
}

func TestFormatVal(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:" min = 0 ; max = 10 ; bla= "`
	}{
		field: -1,
	}) {
		t.Errorf("validators with spaces does not validate")
	}

	if nil != Validate(struct {
		field int `validate:" min = 0 ; max = 10 ; bla = "`
	}{
		field: 5,
	}) {
		t.Errorf("validators with spaces does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"1234567890=!@#$%^&*()"`
	}{
		field: 5,
	}) {
		t.Errorf("incorrect validator must not validate")
	}

	if nil != Validate(struct {
		field int `validate:""`
	}{
		field: 5,
	}) {
		t.Errorf("empty validator must not validate")
	}

	if nil == Validate(struct {
		field int `validate:" one_of = 1 , 2 , 3 "`
	}{
		field: 4,
	}) {
		t.Errorf("validators with spaces does not validate")
	}

	if nil != Validate(struct {
		field int `validate:" one_of = 1 , 2 , 3 "`
	}{
		field: 2,
	}) {
		t.Errorf("validators with spaces does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"one_of="`
	}{
		field: 0,
	}) {
		t.Errorf("empty one_of validate should not validate")
	}
}
func TestMinValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"min=0s"`
	}{
		field: -time.Second,
	}) {
		t.Errorf("min validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"min=-1s"`
	}{
		field: -time.Minute,
	}) {
		t.Errorf("min validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"min=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"min=-1s"`
	}{
		field: -time.Millisecond,
	}) {
		t.Errorf("min validator does not validate for time.Duratuon")
	}
}

func TestMaxValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"max=0s"`
	}{
		field: time.Second,
	}) {
		t.Errorf("max validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"max=1s"`
	}{
		field: time.Minute,
	}) {
		t.Errorf("max validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"max=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"max=1s"`
	}{
		field: time.Millisecond,
	}) {
		t.Errorf("max validator does not validate for time.Duratuon")
	}
}

func TestMinValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"min=0"`
	}{
		field: -1,
	}) {
		t.Errorf("min validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"min=0"`
	}{
		field: -1,
	}) {
		t.Errorf("min validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"min=0"`
	}{
		field: -1,
	}) {
		t.Errorf("min validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"min=0"`
	}{
		field: -1,
	}) {
		t.Errorf("min validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"min=0"`
	}{
		field: -1,
	}) {
		t.Errorf("min validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for int64")
	}
}

func TestMaxValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for int64")
	}
}

func TestMinValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"min=0"`
	}{
		field: -1,
	}) {
		t.Errorf("min validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for rune")
	}
}

func TestMaxValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for rune")
	}
}

func TestMinValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"min=10"`
	}{
		field: 9,
	}) {
		t.Errorf("min validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"min=10"`
	}{
		field: 9,
	}) {
		t.Errorf("min validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"min=10"`
	}{
		field: 9,
	}) {
		t.Errorf("min validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"min=10"`
	}{
		field: 9,
	}) {
		t.Errorf("min validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"min=10"`
	}{
		field: 9,
	}) {
		t.Errorf("min validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"min=10"`
	}{
		field: 9,
	}) {
		t.Errorf("min validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"min=10"`
	}{
		field: 10,
	}) {
		t.Errorf("min validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"min=10"`
	}{
		field: 10,
	}) {
		t.Errorf("min validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"min=10"`
	}{
		field: 10,
	}) {
		t.Errorf("min validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"min=10"`
	}{
		field: 10,
	}) {
		t.Errorf("min validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"min=10"`
	}{
		field: 10,
	}) {
		t.Errorf("min validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uintptr `validate:"min=10"`
	}{
		field: 10,
	}) {
		t.Errorf("min validator does not validate for uintptr")
	}
}

func TestMaxValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"max=0"`
	}{
		field: 1,
	}) {
		t.Errorf("max validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uint64 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for uintptr")
	}
}

func TestMinValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"min=0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("min validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"min=0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("min validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"min=0"`
	}{
		field: 0,
	}) {
		t.Errorf("min validator does not validate for flaot64")
	}
}

func TestMaxValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"max=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("max validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"max=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("max validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"max=0"`
	}{
		field: 0,
	}) {
		t.Errorf("max validator does not validate for flaot64")
	}
}

func TestMinValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"min=2"`
	}{
		field: "a",
	}) {
		t.Errorf("min validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"min=2"`
	}{
		field: "ab",
	}) {
		t.Errorf("min validator does not validate for string")
	}
}

func TestMaxValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"max=2"`
	}{
		field: "abc",
	}) {
		t.Errorf("max validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"max=2"`
	}{
		field: "ab",
	}) {
		t.Errorf("max validator does not validate for string")
	}
}

func TestMinValForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"min=2"`
	}{
		field: map[string]string{
			"a": "a",
		},
	}) {
		t.Errorf("min validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"min=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("min validator does not validate for map")
	}
}

func TestMaxForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"max=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	}) {
		t.Errorf("max validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"max=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("max validator does not validate for map")
	}
}

func TestMinValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"min=2"`
	}{
		field: []string{"a"},
	}) {
		t.Errorf("min validator does not validate for string")
	}

	if nil != Validate(struct {
		field []string `validate:"min=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("min validator does not validate for string")
	}
}

func TestMaxValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"max=2"`
	}{
		field: []string{"a", "b", "c"},
	}) {
		t.Errorf("min validator does not validate for string")
	}

	if nil != Validate(struct {
		field []string `validate:"max=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("min validator does not validate for string")
	}
}

func TestEmptyValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"empty=true"`
	}{
		field: "a",
	}) {
		t.Errorf("empty validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"empty=true"`
	}{
		field: "",
	}) {
		t.Errorf("empty validator does not validate for string")
	}

	if nil == Validate(struct {
		field string `validate:"empty=false"`
	}{
		field: "",
	}) {
		t.Errorf("empty validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"empty=false"`
	}{
		field: "a",
	}) {
		t.Errorf("empty validator does not validate for string")
	}
}

func TestEmptyValForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"empty=true"`
	}{
		field: map[string]string{"a": "a"},
	}) {
		t.Errorf("empty validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"empty=true"`
	}{
		field: map[string]string{},
	}) {
		t.Errorf("empty validator does not validate for map")
	}

	if nil == Validate(struct {
		field map[string]string `validate:"empty=false"`
	}{
		field: map[string]string{},
	}) {
		t.Errorf("empty validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"empty=false"`
	}{
		field: map[string]string{"a": "a"},
	}) {
		t.Errorf("empty validator does not validate for map")
	}
}

func TestEmptyValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"empty=true"`
	}{
		field: []string{
			"a",
		},
	}) {
		t.Errorf("empty validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"empty=true"`
	}{
		field: []string{},
	}) {
		t.Errorf("empty validator does not validate for sclie")
	}

	if nil == Validate(struct {
		field []string `validate:"empty=false"`
	}{
		field: []string{},
	}) {
		t.Errorf("empty validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"empty=false"`
	}{
		field: []string{"a"},
	}) {
		t.Errorf("empty validator does not validate for slice")
	}
}

func TestNilValForPtr(t *testing.T) {
	if nil == Validate(struct {
		field *int `validate:"nil=true"`
	}{
		field: new(int),
	}) {
		t.Errorf("nil validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:"nil=true"`
	}{
		field: nil,
	}) {
		t.Errorf("nil validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field *int `validate:"nil=false"`
	}{
		field: nil,
	}) {
		t.Errorf("nil validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:"nil=false"`
	}{
		field: new(int),
	}) {
		t.Errorf("nil validator does not validate for pointer")
	}
}

func TestOneOfValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"one_of=1s,2s,3s"`
	}{
		field: 4 * time.Second,
	}) {
		t.Errorf("one_of validator does not validate for time.Duration")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"one_of=1s,2s,3s"`
	}{
		field: 2 * time.Second,
	}) {
		t.Errorf("one_of validator does not validate for time.Duration")
	}
}

func TestOneOfValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for int")
	}

	if nil != Validate(struct {
		field int `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int8 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int16 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int32 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int64 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for int64")
	}
}

func TestOneOfValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint8 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint16 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint32 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"one_of=1,2,3"`
	}{
		field: 4,
	}) {
		t.Errorf("one_of validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uint64 `validate:"one_of=1,2,3"`
	}{
		field: 2,
	}) {
		t.Errorf("one_of validator does not validate for uint64")
	}
}

func TestOneOfValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"one_of=1.0,2.0,3.0"`
	}{
		field: 4.0,
	}) {
		t.Errorf("one_of validator does not validate for float32")
	}

	if nil != Validate(struct {
		field float32 `validate:"one_of=1.0,2.0,3.0"`
	}{
		field: 2.0,
	}) {
		t.Errorf("one_of validator does not validate for float32")
	}

	if nil == Validate(struct {
		field float64 `validate:"one_of=1.0,2.0,3.0"`
	}{
		field: 4.0,
	}) {
		t.Errorf("one_of validator does not validate for float64")
	}

	if nil != Validate(struct {
		field float64 `validate:"one_of=1.0,2.0,3.0"`
	}{
		field: 2.0,
	}) {
		t.Errorf("one_of validator does not validate for float64")
	}
}

func TestOneOfValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"one_of=one,two,three"`
	}{
		field: "four",
	}) {
		t.Errorf("one_of validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"one_of=one,two,three"`
	}{
		field: "two",
	}) {
		t.Errorf("one_of validator does not validate for string")
	}
}

func TestDeepValsForMapKeys(t *testing.T) {
	s := " "

	if nil == Validate(struct {
		field map[int]int `validate:"[min=0]"`
	}{
		field: map[int]int{0: 0, -1: 0},
	}) {
		t.Errorf("[min] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"[min=0]"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf("[min] validator does not validate map key")
	}

	if nil == Validate(struct {
		field map[int]int `validate:"[max=0]"`
	}{
		field: map[int]int{0: 0, 1: 0},
	}) {
		t.Errorf("[max] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"[max=0]"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf("[max] validator does not validate for map key")
	}

	if nil == Validate(struct {
		field map[string]int `validate:"[empty=true]"`
	}{
		field: map[string]int{
			" ": 0,
		},
	}) {
		t.Errorf("[empty] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[string]int `validate:"[empty=true]"`
	}{
		field: map[string]int{
			"": 0,
		},
	}) {
		t.Errorf("[empty] validator does not validate for map key")
	}

	if nil == Validate(struct {
		field map[string]int `validate:"[empty=false]"`
	}{
		field: map[string]int{
			"": 0,
		},
	}) {
		t.Errorf("[empty] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[string]int `validate:"[empty=false]"`
	}{
		field: map[string]int{
			" ": 0,
		},
	}) {
		t.Errorf("[empty] validator does not validate for map key")
	}

	if nil == Validate(struct {
		field map[*string]int `validate:"[nil=true]"`
	}{
		field: map[*string]int{
			&s: 0,
		},
	}) {
		t.Errorf("[nil] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[*string]int `validate:"[nil=true]"`
	}{
		field: map[*string]int{
			nil: 0,
		},
	}) {
		t.Errorf("[nil] validator does not validate for map key")
	}

	if nil == Validate(struct {
		field map[*string]int `validate:"[nil=false]"`
	}{
		field: map[*string]int{
			nil: 0,
		},
	}) {
		t.Errorf("[nil] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[*string]int `validate:"[nil=false]"`
	}{
		field: map[*string]int{
			&s: 0,
		},
	}) {
		t.Errorf("[nil] validator does not validate for map key")
	}

	if nil == Validate(struct {
		field map[int]int `validate:"[one_of=1,2,3]"`
	}{
		field: map[int]int{
			4: 0,
		},
	}) {
		t.Errorf("[one_of] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"[one_of=1,2,3]"`
	}{
		field: map[int]int{
			1: 0,
			2: 0,
			3: 0,
		},
	}) {
		t.Errorf("[one_of] validator does not validate for map key")
	}
}

func TestDeepValsForMapValues(t *testing.T) {
	s := " "

	if nil == Validate(struct {
		field map[int]int `validate:"> min=0"`
	}{
		field: map[int]int{0: 0, 1: -1},
	}) {
		t.Errorf(">min validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"> min=0"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf(">min validator does not validate map values")
	}

	if nil == Validate(struct {
		field map[int]int `validate:"> max=0"`
	}{
		field: map[int]int{0: 0, -1: 1},
	}) {
		t.Errorf(">max validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"> max=0"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf(">max validator does not validate for map values")
	}

	if nil == Validate(struct {
		field map[int]string `validate:"> empty=true"`
	}{
		field: map[int]string{
			0: " ",
		},
	}) {
		t.Errorf(">empty validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]string `validate:"> empty=true"`
	}{
		field: map[int]string{
			0: "",
		},
	}) {
		t.Errorf(">empty validator does not validate for map values")
	}

	if nil == Validate(struct {
		field map[int]string `validate:"> empty=false"`
	}{
		field: map[int]string{
			0: "",
		},
	}) {
		t.Errorf(">empty validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]string `validate:"> empty=false"`
	}{
		field: map[int]string{
			0: " ",
		},
	}) {
		t.Errorf(">empty validator does not validate for map values")
	}

	if nil == Validate(struct {
		field map[int]*string `validate:"> nil=true"`
	}{
		field: map[int]*string{
			0: &s,
		},
	}) {
		t.Errorf(">nil validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]*string `validate:"> nil=true"`
	}{
		field: map[int]*string{
			0: nil,
		},
	}) {
		t.Errorf(">nil validator does not validate for map values")
	}

	if nil == Validate(struct {
		field map[int]*string `validate:"> nil=false"`
	}{
		field: map[int]*string{
			0: nil,
		},
	}) {
		t.Errorf(">nil validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]*string `validate:"> nil=false]"`
	}{
		field: map[int]*string{
			0: &s,
		},
	}) {
		t.Errorf(">nil validator does not validate for map values")
	}

	if nil == Validate(struct {
		field map[int]int `validate:"> one_of=1,2,3"`
	}{
		field: map[int]int{
			0: 4,
		},
	}) {
		t.Errorf(">one_of validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"> one_of=1,2,3"`
	}{
		field: map[int]int{
			-1: 1,
			-2: 2,
			-3: 3,
		},
	}) {
		t.Errorf(">one_of validator does not validate for map values")
	}

}

func TestDeepValsForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []int `validate:">min=0"`
	}{
		field: []int{0, -1},
	}) {
		t.Errorf(">min validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []int `validate:">min=0"`
	}{
		field: []int{0, 0},
	}) {
		t.Errorf(">min validator does not validate for slice")
	}

	if nil == Validate(struct {
		field []int `validate:">max=0"`
	}{
		field: []int{0, 1},
	}) {
		t.Errorf(">max validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []int `validate:">max=0"`
	}{
		field: []int{0, 0},
	}) {
		t.Errorf(">max validator does not validate for slice")
	}

	if nil == Validate(struct {
		field [][]int `validate:">empty=true"`
	}{
		field: [][]int{
			[]int{0},
		},
	}) {
		t.Errorf(">empty validator does not validate for slice")
	}

	if nil != Validate(struct {
		field [][]int `validate:">empty=true"`
	}{
		field: [][]int{
			[]int{},
		},
	}) {
		t.Errorf(">empty validator does not validate for slice")
	}

	if nil == Validate(struct {
		field [][]int `validate:">empty=false"`
	}{
		field: [][]int{
			[]int{},
		},
	}) {
		t.Errorf(">empty validator does not validate for slice")
	}

	if nil != Validate(struct {
		field [][]int `validate:">empty=false"`
	}{
		field: [][]int{
			[]int{0},
		},
	}) {
		t.Errorf(">empty validator does not validate for slice")
	}

	if nil == Validate(struct {
		field []*int `validate:">nil=true"`
	}{
		field: []*int{
			new(int),
		},
	}) {
		t.Errorf(">nil validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []*int `validate:">nil=true"`
	}{
		field: []*int{nil},
	}) {
		t.Errorf(">nil validator does not validate for slice")
	}

	if nil == Validate(struct {
		field []*int `validate:">nil=false"`
	}{
		field: []*int{nil},
	}) {
		t.Errorf(">nil validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []*int `validate:">nil=false"`
	}{
		field: []*int{new(int)},
	}) {
		t.Errorf(">nil validator does not validate for slice")
	}

	if nil == Validate(struct {
		field []int `validate:">one_of=1,2,3"`
	}{
		field: []int{4},
	}) {
		t.Errorf(">one_of validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []int `validate:">one_of=1,2,3"`
	}{
		field: []int{1, 2, 3},
	}) {
		t.Errorf(">one_of validator does not validate for slice")
	}
}

func TestDeepValsForPtr(t *testing.T) {
	minusOne := -1
	zero := 0
	one := 1
	four := 4
	empty := ""
	notEmpty := "a"
	onePtr := &one
	var nilPtr *int

	if nil == Validate(struct {
		field *int `validate:">min=0"`
	}{
		field: &minusOne,
	}) {
		t.Errorf(">min validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:">min=0"`
	}{
		field: &zero,
	}) {
		t.Errorf(">min validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field *int `validate:">max=0"`
	}{
		field: &one,
	}) {
		t.Errorf(">max validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:">max=0"`
	}{
		field: &zero,
	}) {
		t.Errorf(">max validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field *string `validate:">empty=true"`
	}{
		field: &notEmpty,
	}) {
		t.Errorf(">empty validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *string `validate:">empty=true"`
	}{
		field: &empty,
	}) {
		t.Errorf(">empty validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field *string `validate:">empty=false"`
	}{
		field: &empty,
	}) {
		t.Errorf(">empty validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *string `validate:">empty=false"`
	}{
		field: &notEmpty,
	}) {
		t.Errorf(">empty validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field **int `validate:">nil=true"`
	}{
		field: &onePtr,
	}) {
		t.Errorf(">nil validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field **int `validate:">nil=true"`
	}{
		field: &nilPtr,
	}) {
		t.Errorf(">nil validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field **int `validate:">nil=false"`
	}{
		field: &nilPtr,
	}) {
		t.Errorf(">nil validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field **int `validate:">nil=false"`
	}{
		field: &onePtr,
	}) {
		t.Errorf(">nil validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field *int `validate:">one_of=1,2,3"`
	}{
		field: &four,
	}) {
		t.Errorf(">one_of validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:">one_of=1,2,3"`
	}{
		field: &one,
	}) {
		t.Errorf(">one_of validator does not validate for pointer")
	}
}

func TestDeepVal(t *testing.T) {
	// Should not validate one level deep

	if nil != Validate(struct {
		field []int `validate:"min=0"`
	}{
		field: []int{-1},
	}) {
		t.Errorf("min validator validates one level deep for slice")
	}

	one := 1
	if nil != Validate(struct {
		field *int `validate:"max=0"`
	}{
		field: &one,
	}) {
		t.Errorf("max validator validates one level deep for pointer")
	}
}
