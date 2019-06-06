package validate

import (
	"errors"
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
	var validatorsOr [][]validator

	validatorsOr = parseValidators("")
	temp := make([][]validator, 0)
	if !reflect.DeepEqual(validatorsOr, temp) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	validatorsOr = parseValidators("&|&,&")
	if !reflect.DeepEqual(validatorsOr, temp) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	validatorsOr = parseValidators("val_a=a")
	if !reflect.DeepEqual(validatorsOr, [][]validator{
		[]validator{
			validator{
				ValidatorType("val_a"),
				"a",
			},
		},
	}) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	validatorsOr = parseValidators("  val  |val_a=a|val_1 = 1  |  val_b = b , c_d_ , 1.0  |VAL = V A L U E ¶  ")
	if !reflect.DeepEqual(validatorsOr, [][]validator{
		[]validator{
			validator{
				ValidatorType("val"),
				"",
			},
		},
		[]validator{
			validator{
				ValidatorType("val_a"),
				"a",
			},
		},
		[]validator{
			validator{
				ValidatorType("val_1"),
				"1",
			},
		},
		[]validator{
			validator{
				ValidatorType("val_b"),
				"b , c_d_ , 1.0",
			},
		},
		[]validator{
			validator{
				ValidatorType("VAL"),
				"V A L U E ¶",
			},
		},
	}) {
		t.Errorf("parseValidators incorrectly parses validators")
	}

	validatorsOr = parseValidators("  val  &val_a=a|val_1 = 1  &  val_b = b , c_d_ , 1.0  &VAL = V A L U E ¶  ")
	if !reflect.DeepEqual(validatorsOr, [][]validator{
		[]validator{
			validator{
				ValidatorType("val"),
				"",
			},
			validator{
				ValidatorType("val_a"),
				"a",
			},
		},
		[]validator{
			validator{
				ValidatorType("val_1"),
				"1",
			},
			validator{
				ValidatorType("val_b"),
				"b , c_d_ , 1.0",
			},
			validator{
				ValidatorType("VAL"),
				"V A L U E ¶",
			},
		},
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
	type St struct {
		field int
	}
	st := St{
		field: 1,
	}
	if nil != Validate(st) {
		t.Errorf("validate does not validate struct type")
	}
	if nil != Validate(&st) {
		t.Errorf("validate does not validate struct pointer type")
	}
	if nil != Validate(map[int]St{
		0: st,
	}) {
		t.Errorf("validate does not validate map of struct type")
	}
	if nil != Validate(map[St]int{
		st: 0,
	}) {
		t.Errorf("validate does not validate map of struct type")
	}
	if nil != Validate([]St{
		st,
	}) {
		t.Errorf("validate does not validate slice of struct type")
	}
	if nil != Validate([1]St{
		st,
	}) {
		t.Errorf("validate does not validate slice of struct type")
	}

	type StFail struct {
		field int `validate:"lte=0"`
	}
	stFail := StFail{
		field: 1,
	}
	if nil == Validate(stFail) {
		t.Errorf("validate does not validate struct type")
	}
	if nil == Validate(&stFail) {
		t.Errorf("validate does not validate struct pointer type")
	}
	if nil == Validate(map[int]StFail{
		0: stFail,
	}) {
		t.Errorf("validate does not validate map of struct type")
	}
	if nil == Validate(map[StFail]int{
		stFail: 0,
	}) {
		t.Errorf("validate does not validate map of struct type")
	}
	if nil == Validate([]StFail{
		stFail,
	}) {
		t.Errorf("validate does not validate slice of struct type")
	}
	if nil == Validate([1]StFail{
		stFail,
	}) {
		t.Errorf("validate does not validate slice of struct type")
	}

	type StAnotherFail struct {
		a     int
		b     int
		field int `validate:"lte=0"`
		c     int
		d     int
	}
	stAnotherFail := StAnotherFail{
		field: 1,
	}
	if nil == Validate(stAnotherFail) {
		t.Errorf("validate does not validate struct type")
	}
	if nil == Validate(&stAnotherFail) {
		t.Errorf("validate does not validate struct pointer type")
	}
	if nil == Validate(map[int]StAnotherFail{
		0: stAnotherFail,
	}) {
		t.Errorf("validate does not validate map of struct type")
	}
	if nil == Validate(map[StAnotherFail]int{
		stAnotherFail: 0,
	}) {
		t.Errorf("validate does not validate map of struct type")
	}
	if nil == Validate([]StAnotherFail{
		stAnotherFail,
	}) {
		t.Errorf("validate does not validate slice of struct type")
	}
	if nil == Validate([1]StAnotherFail{
		stAnotherFail,
	}) {
		t.Errorf("validate does not validate slice of struct type")
	}
}

func TestErrors(t *testing.T) {
	err := Validate([]struct {
		field map[time.Duration]int `validate:"gte=2 [eq=0s]"`
	}{{
		field: map[time.Duration]int{-time.Second: 1},
	}})

	switch err.(type) {
	case ErrorValidation:
	default:
		t.Errorf("error of the wrong type")
	}
}

type StCustomValidator struct {
	field        int
	anotherField int `validate:"eq=0"`
}

func (st StCustomValidator) Validate() error {
	if st.field <= 0 {
		return errors.New("field should be positive")
	}

	return nil
}

type StCustomValidator2 struct {
	field        int
	anotherField int `validate:"eq=0"`
}

func (st *StCustomValidator2) Validate() error {
	if st.field <= 0 {
		return errors.New("field should be positive")
	}

	return nil
}

type IntCustomValidator int

func (i IntCustomValidator) Validate() error {
	if i <= 0 {
		return errors.New("field should be positive")
	}

	return nil
}

func TestCustomValidator(t *testing.T) {
	// Test a custom validtor wiht a value reciever by value, then by reference
	if nil != Validate(StCustomValidator{
		field: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(StCustomValidator{
		field: 0,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(&StCustomValidator{
		field: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&StCustomValidator{
		field: 0,
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test a custom validtor with a pointer reciever by value, then by reference
	if nil != Validate(StCustomValidator2{
		field: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(StCustomValidator2{
		field: 0,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(&StCustomValidator2{
		field: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&StCustomValidator2{
		field: 0,
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test an embedded validator together with a custom validtor with a value reciever by value, then by reference
	if nil == Validate(StCustomValidator{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(StCustomValidator{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&StCustomValidator{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&StCustomValidator{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test an embedded validator together with a custom validtor with a pointer reciever by value, then by reference
	if nil == Validate(StCustomValidator2{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(StCustomValidator2{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&StCustomValidator2{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&StCustomValidator2{
		field:        1,
		anotherField: 1,
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test a custom validtor of a subtruct (defined in exported field) with a value reciever by value, then by reference
	if nil != Validate(struct {
		Field StCustomValidator
	}{
		Field: StCustomValidator{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(struct {
		Field StCustomValidator
	}{
		Field: StCustomValidator{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		Field *StCustomValidator
	}{
		Field: &StCustomValidator{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(struct {
		Field *StCustomValidator
	}{
		Field: &StCustomValidator{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test a custom validtor of a subtruct (defined in exported field)  with a pointer reciever by value, then by reference
	if nil != Validate(struct {
		Field StCustomValidator2
	}{
		Field: StCustomValidator2{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(struct {
		Field StCustomValidator2
	}{
		Field: StCustomValidator2{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		Field *StCustomValidator2
	}{
		Field: &StCustomValidator2{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(struct {
		Field *StCustomValidator2
	}{
		Field: &StCustomValidator2{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test a custom validtor of a subtruct (defined in unexported field) with a value reciever by value, then by reference
	// This will fail, but should not panic
	if nil != Validate(struct {
		field StCustomValidator
	}{
		field: StCustomValidator{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		field StCustomValidator
	}{
		field: StCustomValidator{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		field *StCustomValidator
	}{
		field: &StCustomValidator{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		field *StCustomValidator
	}{
		field: &StCustomValidator{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test a custom validtor of a subtruct (defined in unexported field) with a pointer reciever by value, then by reference
	// This will fail, but should not panic
	if nil != Validate(struct {
		field StCustomValidator2
	}{
		field: StCustomValidator2{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		field StCustomValidator2
	}{
		field: StCustomValidator2{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		field *StCustomValidator2
	}{
		field: &StCustomValidator2{
			field: 1,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(struct {
		field *StCustomValidator2
	}{
		field: &StCustomValidator2{
			field: 0,
		},
	}) {
		t.Errorf("custom validator does not validate")
	}

	// Test a custom validtor of an artbitrary type
	one := IntCustomValidator(1)
	zero := IntCustomValidator(0)

	if nil != Validate(one) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(zero) {
		t.Errorf("custom validator does not validate")
	}

	if nil != Validate(&one) {
		t.Errorf("custom validator does not validate")
	}

	if nil == Validate(&zero) {
		t.Errorf("custom validator does not validate")
	}
}

func TestAndVal(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"gte=0&lte=10"`
	}{
		field: -1,
	}) {
		t.Errorf("& operator does not validate")
	}

	if nil == Validate(struct {
		field int `validate:"gte=0&lte=10"`
	}{
		field: 11,
	}) {
		t.Errorf("& operator does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0&lte=10"`
	}{
		field: 5,
	}) {
		t.Errorf("& operator does not validate")
	}

	if nil == Validate(struct {
		field int `validate:"gte=1&lte=-1"`
	}{
		field: 0,
	}) {
		t.Errorf("& operator does not validate")
	}
}

func TestOrVal(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"lte=0|gte=10"`
	}{
		field: 5,
	}) {
		t.Errorf("| operator does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"lte=0|gte=10"`
	}{
		field: -1,
	}) {
		t.Errorf("| operator does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"lte=0|gte=10"`
	}{
		field: 11,
	}) {
		t.Errorf("| operator does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0|lte=10"`
	}{
		field: 5,
	}) {
		t.Errorf("| operator does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0|lte=10"`
	}{
		field: -1,
	}) {
		t.Errorf("| operator does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0|lte=10"`
	}{
		field: 11,
	}) {
		t.Errorf("| operator does not validate")
	}
}

func TestAndOrVal(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"gte=0 & lte=5 | gte=10 & lte=15"`
	}{
		field: -1,
	}) {
		t.Errorf("& and | operators does not validate")
	}

	if nil == Validate(struct {
		field int `validate:"gte=0 & lte=5 | gte=10 & lte=15"`
	}{
		field: 6,
	}) {
		t.Errorf("& and | operators does not validate")
	}

	if nil == Validate(struct {
		field int `validate:"gte=0 & lte=5 | gte=10 & lte=15"`
	}{
		field: 16,
	}) {
		t.Errorf("& and | operators does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0 & lte=5 | gte=10 & lte=15"`
	}{
		field: 1,
	}) {
		t.Errorf("& and | operators does not validate")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0 & lte=5 | gte=10 & lte=15"`
	}{
		field: 11,
	}) {
		t.Errorf("& and | operators does not validate")
	}
}

func TestFormatVal(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:" gte = 0 & lte = 10 & bla= "`
	}{
		field: -1,
	}) {
		t.Errorf("validators with spaces does not validate")
	}

	if nil != Validate(struct {
		field int `validate:" gte = 0 & lte = 10 & bla = "`
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

func TestEqValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"eq=0s"`
	}{
		field: -time.Second,
	}) {
		t.Errorf("eq validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"eq=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for time.Duratuon")
	}
}

func TestEqValForInt(t *testing.T) {
	if nil != Validate(struct {
		field int `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for int64")
	}

	if nil == Validate(struct {
		field int `validate:"eq=0"`
	}{
		field: 1,
	}) {
		t.Errorf("eq validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"eq=0"`
	}{
		field: 1,
	}) {
		t.Errorf("eq validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"eq=0"`
	}{
		field: 1,
	}) {
		t.Errorf("eq validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"eq=0"`
	}{
		field: 1,
	}) {
		t.Errorf("eq validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"eq=0"`
	}{
		field: 1,
	}) {
		t.Errorf("eq validator does not validate for int64")
	}
}

func TestEqValForRune(t *testing.T) {
	if nil != Validate(struct {
		field rune `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for rune")
	}

	if nil == Validate(struct {
		field rune `validate:"eq=0"`
	}{
		field: 1,
	}) {
		t.Errorf("eq validator does not validate for rune")
	}
}

func TestEqValForUint(t *testing.T) {
	if nil != Validate(struct {
		field uint `validate:"eq=10"`
	}{
		field: 10,
	}) {
		t.Errorf("eq validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"eq=10"`
	}{
		field: 10,
	}) {
		t.Errorf("eq validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"eq=10"`
	}{
		field: 10,
	}) {
		t.Errorf("eq validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"eq=10"`
	}{
		field: 10,
	}) {
		t.Errorf("eq validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"eq=10"`
	}{
		field: 10,
	}) {
		t.Errorf("eq validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uintptr `validate:"eq=10"`
	}{
		field: 10,
	}) {
		t.Errorf("eq validator does not validate for uintptr")
	}

	if nil == Validate(struct {
		field uint `validate:"eq=10"`
	}{
		field: 11,
	}) {
		t.Errorf("eq validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"eq=10"`
	}{
		field: 11,
	}) {
		t.Errorf("eq validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"eq=10"`
	}{
		field: 11,
	}) {
		t.Errorf("eq validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"eq=10"`
	}{
		field: 11,
	}) {
		t.Errorf("eq validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"eq=10"`
	}{
		field: 11,
	}) {
		t.Errorf("eq validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"eq=10"`
	}{
		field: 11,
	}) {
		t.Errorf("eq validator does not validate for uintptr")
	}
}

func TestEqValForFloat(t *testing.T) {
	if nil != Validate(struct {
		field float32 `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"eq=0"`
	}{
		field: 0,
	}) {
		t.Errorf("eq validator does not validate for flaot64")
	}

	if nil == Validate(struct {
		field float32 `validate:"eq=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("eq validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"eq=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("eq validator does not validate for flaot64")
	}
}

func TestEqValForString(t *testing.T) {
	if nil != Validate(struct {
		field string `validate:"eq=2"`
	}{
		field: "aa",
	}) {
		t.Errorf("eq validator does not validate for string")
	}

	if nil == Validate(struct {
		field string `validate:"eq=2"`
	}{
		field: "abc",
	}) {
		t.Errorf("eq validator does not validate for string")
	}
}

func TestEqValForMap(t *testing.T) {
	if nil != Validate(struct {
		field map[string]string `validate:"eq=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("eq validator does not validate for map")
	}

	if nil == Validate(struct {
		field map[string]string `validate:"eq=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	}) {
		t.Errorf("eq validator does not validate for map")
	}
}

func TestEqValForSlice(t *testing.T) {
	if nil != Validate(struct {
		field []string `validate:"eq=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("eq validator does not validate for slice")
	}

	if nil == Validate(struct {
		field []string `validate:"eq=2"`
	}{
		field: []string{"a", "b", "c"},
	}) {
		t.Errorf("eq validator does not validate for slice")
	}
}

func TestEqValForArray(t *testing.T) {
	if nil != Validate(struct {
		field [2]string `validate:"eq=2"`
	}{
		field: [2]string{"a", "b"},
	}) {
		t.Errorf("eq validator does not validate for string")
	}

	if nil == Validate(struct {
		field [3]string `validate:"eq=2"`
	}{
		field: [3]string{"a", "b", "c"},
	}) {
		t.Errorf("eq validator does not validate for string")
	}
}

func TestNeValForDuration(t *testing.T) {
	if nil != Validate(struct {
		field time.Duration `validate:"ne=0s"`
	}{
		field: -time.Second,
	}) {
		t.Errorf("ne validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"ne=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for time.Duratuon")
	}
}

func TestNeValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"ne=0"`
	}{
		field: 1,
	}) {
		t.Errorf("ne validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"ne=0"`
	}{
		field: 1,
	}) {
		t.Errorf("ne validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"ne=0"`
	}{
		field: 1,
	}) {
		t.Errorf("ne validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"ne=0"`
	}{
		field: 1,
	}) {
		t.Errorf("ne validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"ne=0"`
	}{
		field: 1,
	}) {
		t.Errorf("ne validator does not validate for int64")
	}
}

func TestNeValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"ne=0"`
	}{
		field: 1,
	}) {
		t.Errorf("ne validator does not validate for rune")
	}
}

func TestNeValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"ne=10"`
	}{
		field: 10,
	}) {
		t.Errorf("ne validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"ne=10"`
	}{
		field: 10,
	}) {
		t.Errorf("ne validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"ne=10"`
	}{
		field: 10,
	}) {
		t.Errorf("ne validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"ne=10"`
	}{
		field: 10,
	}) {
		t.Errorf("ne validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"ne=10"`
	}{
		field: 10,
	}) {
		t.Errorf("ne validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"ne=10"`
	}{
		field: 10,
	}) {
		t.Errorf("ne validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"ne=10"`
	}{
		field: 11,
	}) {
		t.Errorf("ne validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"ne=10"`
	}{
		field: 11,
	}) {
		t.Errorf("ne validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"ne=10"`
	}{
		field: 11,
	}) {
		t.Errorf("ne validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"ne=10"`
	}{
		field: 11,
	}) {
		t.Errorf("ne validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"ne=10"`
	}{
		field: 11,
	}) {
		t.Errorf("ne validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uintptr `validate:"ne=10"`
	}{
		field: 11,
	}) {
		t.Errorf("ne validator does not validate for uintptr")
	}
}

func TestNeValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"ne=0"`
	}{
		field: 0,
	}) {
		t.Errorf("ne validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"ne=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("ne validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"ne=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("ne validator does not validate for flaot64")
	}
}

func TestNeValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"ne=2"`
	}{
		field: "aa",
	}) {
		t.Errorf("ne validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"ne=2"`
	}{
		field: "abc",
	}) {
		t.Errorf("ne validator does not validate for string")
	}
}

func TestNeValForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"ne=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("ne validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"ne=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	}) {
		t.Errorf("ne validator does not validate for map")
	}
}

func TestNeValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"ne=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("ne validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"ne=2"`
	}{
		field: []string{"a", "b", "c"},
	}) {
		t.Errorf("ne validator does not validate for slice")
	}
}

func TestNeValForArray(t *testing.T) {
	if nil == Validate(struct {
		field [2]string `validate:"ne=2"`
	}{
		field: [2]string{"a", "b"},
	}) {
		t.Errorf("ne validator does not validate for string")
	}

	if nil != Validate(struct {
		field [3]string `validate:"ne=2"`
	}{
		field: [3]string{"a", "b", "c"},
	}) {
		t.Errorf("ne validator does not validate for string")
	}
}

func TestGtValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"gt=0s"`
	}{
		field: -time.Second,
	}) {
		t.Errorf("gt validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"gt=-1s"`
	}{
		field: -time.Minute,
	}) {
		t.Errorf("gt validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"gt=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"gt=-1s"`
	}{
		field: -time.Millisecond,
	}) {
		t.Errorf("gt validator does not validate for time.Duratuon")
	}
}

func TestLtValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"lt=0s"`
	}{
		field: time.Second,
	}) {
		t.Errorf("lt validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"lt=1s"`
	}{
		field: time.Minute,
	}) {
		t.Errorf("lt validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"lt=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"lt=1s"`
	}{
		field: time.Millisecond,
	}) {
		t.Errorf("lt validator does not validate for time.Duratuon")
	}
}

func TestGtValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"gt=0"`
	}{
		field: 1,
	}) {
		t.Errorf("gt validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"gt=0"`
	}{
		field: 1,
	}) {
		t.Errorf("gt validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"gt=0"`
	}{
		field: 1,
	}) {
		t.Errorf("gt validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"gt=0"`
	}{
		field: 1,
	}) {
		t.Errorf("gt validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"gt=0"`
	}{
		field: 1,
	}) {
		t.Errorf("gt validator does not validate for int64")
	}
}

func TestLtValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"lt=0"`
	}{
		field: -1,
	}) {
		t.Errorf("lt validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"lt=0"`
	}{
		field: -1,
	}) {
		t.Errorf("lt validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"lt=0"`
	}{
		field: -1,
	}) {
		t.Errorf("lt validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"lt=0"`
	}{
		field: -1,
	}) {
		t.Errorf("lt validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"lt=0"`
	}{
		field: -1,
	}) {
		t.Errorf("lt validator does not validate for int64")
	}
}

func TestGtValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"gt=0"`
	}{
		field: 1,
	}) {
		t.Errorf("gt validator does not validate for rune")
	}
}

func TestLtValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"lt=0"`
	}{
		field: -1,
	}) {
		t.Errorf("lt validator does not validate for rune")
	}
}

func TestGtValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"gt=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gt validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"gt=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gt validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"gt=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gt validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"gt=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gt validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"gt=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gt validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"gt=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gt validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"gt=10"`
	}{
		field: 11,
	}) {
		t.Errorf("gt validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"gt=10"`
	}{
		field: 11,
	}) {
		t.Errorf("gt validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"gt=10"`
	}{
		field: 11,
	}) {
		t.Errorf("gt validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"gt=10"`
	}{
		field: 11,
	}) {
		t.Errorf("gt validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"gt=10"`
	}{
		field: 11,
	}) {
		t.Errorf("gt validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uintptr `validate:"gt=10"`
	}{
		field: 11,
	}) {
		t.Errorf("gt validator does not validate for uintptr")
	}
}

func TestLtValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"lt=1"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"lt=1"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"lt=1"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"lt=1"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"lt=1"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uint64 `validate:"lt=1"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for uintptr")
	}
}

func TestGtValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"gt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gt validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"gt=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("gt validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"gt=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("gt validator does not validate for flaot64")
	}
}

func TestLtValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"lt=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lt validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"lt=0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("lt validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"lt=0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("lt validator does not validate for flaot64")
	}
}

func TestGtValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"gt=2"`
	}{
		field: "aa",
	}) {
		t.Errorf("gt validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"gt=2"`
	}{
		field: "abc",
	}) {
		t.Errorf("gt validator does not validate for string")
	}
}

func TestLtValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"lt=2"`
	}{
		field: "ab",
	}) {
		t.Errorf("lt validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"lt=2"`
	}{
		field: "a",
	}) {
		t.Errorf("lt validator does not validate for string")
	}
}

func TestGtValForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"gt=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("gt validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"gt=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	}) {
		t.Errorf("gt validator does not validate for map")
	}
}

func TestLtForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"lt=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("lt validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"lt=2"`
	}{
		field: map[string]string{
			"a": "a",
		},
	}) {
		t.Errorf("lt validator does not validate for map")
	}
}

func TestGtValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"gt=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("gt validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"gt=2"`
	}{
		field: []string{"a", "b", "c"},
	}) {
		t.Errorf("gt validator does not validate for slice")
	}
}

func TestLtValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"lt=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("gt validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"lt=2"`
	}{
		field: []string{"a"},
	}) {
		t.Errorf("gt validator does not validate for slice")
	}
}

func TestGtValForArray(t *testing.T) {
	if nil == Validate(struct {
		field [2]string `validate:"gt=2"`
	}{
		field: [2]string{"a", "b"},
	}) {
		t.Errorf("gt validator does not validate for string")
	}

	if nil != Validate(struct {
		field [3]string `validate:"gt=2"`
	}{
		field: [3]string{"a", "b", "c"},
	}) {
		t.Errorf("gt validator does not validate for string")
	}
}

func TestLtValForArray(t *testing.T) {
	if nil == Validate(struct {
		field [2]string `validate:"lt=2"`
	}{
		field: [2]string{"a", "b"},
	}) {
		t.Errorf("gt validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1]string `validate:"lt=2"`
	}{
		field: [1]string{"a"},
	}) {
		t.Errorf("gt validator does not validate for array")
	}
}

func TestGteValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"gte=0s"`
	}{
		field: -time.Second,
	}) {
		t.Errorf("gte validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"gte=-1s"`
	}{
		field: -time.Minute,
	}) {
		t.Errorf("gte validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"gte=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"gte=-1s"`
	}{
		field: -time.Millisecond,
	}) {
		t.Errorf("gte validator does not validate for time.Duratuon")
	}
}

func TestLteValForDuration(t *testing.T) {
	if nil == Validate(struct {
		field time.Duration `validate:"lte=0s"`
	}{
		field: time.Second,
	}) {
		t.Errorf("lte validator does not validate for time.Duratuon")
	}

	if nil == Validate(struct {
		field time.Duration `validate:"lte=1s"`
	}{
		field: time.Minute,
	}) {
		t.Errorf("lte validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"lte=0s"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for time.Duratuon")
	}

	if nil != Validate(struct {
		field time.Duration `validate:"lte=1s"`
	}{
		field: time.Millisecond,
	}) {
		t.Errorf("lte validator does not validate for time.Duratuon")
	}
}

func TestGteValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"gte=0"`
	}{
		field: -1,
	}) {
		t.Errorf("gte validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"gte=0"`
	}{
		field: -1,
	}) {
		t.Errorf("gte validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"gte=0"`
	}{
		field: -1,
	}) {
		t.Errorf("gte validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"gte=0"`
	}{
		field: -1,
	}) {
		t.Errorf("gte validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"gte=0"`
	}{
		field: -1,
	}) {
		t.Errorf("gte validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for int64")
	}
}

func TestLteValForInt(t *testing.T) {
	if nil == Validate(struct {
		field int `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for int")
	}

	if nil == Validate(struct {
		field int8 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for int8")
	}

	if nil == Validate(struct {
		field int16 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for int16")
	}

	if nil == Validate(struct {
		field int32 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for int32")
	}

	if nil == Validate(struct {
		field int64 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for int64")
	}

	if nil != Validate(struct {
		field int `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for int")
	}

	if nil != Validate(struct {
		field int8 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for int8")
	}

	if nil != Validate(struct {
		field int16 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for int16")
	}

	if nil != Validate(struct {
		field int32 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for int32")
	}

	if nil != Validate(struct {
		field int64 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for int64")
	}
}

func TestGteValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"gte=0"`
	}{
		field: -1,
	}) {
		t.Errorf("gte validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for rune")
	}
}

func TestLteValForRune(t *testing.T) {
	if nil == Validate(struct {
		field rune `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for rune")
	}

	if nil != Validate(struct {
		field rune `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for rune")
	}
}

func TestGteValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"gte=10"`
	}{
		field: 9,
	}) {
		t.Errorf("gte validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"gte=10"`
	}{
		field: 9,
	}) {
		t.Errorf("gte validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"gte=10"`
	}{
		field: 9,
	}) {
		t.Errorf("gte validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"gte=10"`
	}{
		field: 9,
	}) {
		t.Errorf("gte validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"gte=10"`
	}{
		field: 9,
	}) {
		t.Errorf("gte validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"gte=10"`
	}{
		field: 9,
	}) {
		t.Errorf("gte validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"gte=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gte validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"gte=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gte validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"gte=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gte validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"gte=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gte validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"gte=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gte validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uintptr `validate:"gte=10"`
	}{
		field: 10,
	}) {
		t.Errorf("gte validator does not validate for uintptr")
	}
}

func TestLteValForUint(t *testing.T) {
	if nil == Validate(struct {
		field uint `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for uint")
	}

	if nil == Validate(struct {
		field uint8 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for uint8")
	}

	if nil == Validate(struct {
		field uint16 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for uint16")
	}

	if nil == Validate(struct {
		field uint32 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for uint32")
	}

	if nil == Validate(struct {
		field uint64 `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for uint64")
	}

	if nil == Validate(struct {
		field uintptr `validate:"lte=0"`
	}{
		field: 1,
	}) {
		t.Errorf("lte validator does not validate for uintptr")
	}

	if nil != Validate(struct {
		field uint `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for uint")
	}

	if nil != Validate(struct {
		field uint8 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for uint8")
	}

	if nil != Validate(struct {
		field uint16 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for uint16")
	}

	if nil != Validate(struct {
		field uint32 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for uint32")
	}

	if nil != Validate(struct {
		field uint64 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for uint64")
	}

	if nil != Validate(struct {
		field uint64 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for uintptr")
	}
}

func TestGteValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"gte=0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("gte validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"gte=0"`
	}{
		field: -0.1,
	}) {
		t.Errorf("gte validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"gte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("gte validator does not validate for flaot64")
	}
}

func TestLteValForFloat(t *testing.T) {
	if nil == Validate(struct {
		field float32 `validate:"lte=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("lte validator does not validate for flaot32")
	}

	if nil == Validate(struct {
		field float64 `validate:"lte=0"`
	}{
		field: 0.1,
	}) {
		t.Errorf("lte validator does not validate for flaot64")
	}

	if nil != Validate(struct {
		field float32 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for flaot32")
	}

	if nil != Validate(struct {
		field float64 `validate:"lte=0"`
	}{
		field: 0,
	}) {
		t.Errorf("lte validator does not validate for flaot64")
	}
}

func TestGteValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"gte=2"`
	}{
		field: "a",
	}) {
		t.Errorf("gte validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"gte=2"`
	}{
		field: "ab",
	}) {
		t.Errorf("gte validator does not validate for string")
	}
}

func TestLteValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"lte=2"`
	}{
		field: "abc",
	}) {
		t.Errorf("lte validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"lte=2"`
	}{
		field: "ab",
	}) {
		t.Errorf("lte validator does not validate for string")
	}
}

func TestGteValForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"gte=2"`
	}{
		field: map[string]string{
			"a": "a",
		},
	}) {
		t.Errorf("gte validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"gte=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("gte validator does not validate for map")
	}
}

func TestLteForMap(t *testing.T) {
	if nil == Validate(struct {
		field map[string]string `validate:"lte=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	}) {
		t.Errorf("lte validator does not validate for map")
	}

	if nil != Validate(struct {
		field map[string]string `validate:"lte=2"`
	}{
		field: map[string]string{
			"a": "a",
			"b": "b",
		},
	}) {
		t.Errorf("lte validator does not validate for map")
	}
}

func TestGteValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"gte=2"`
	}{
		field: []string{"a"},
	}) {
		t.Errorf("gte validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"gte=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("gte validator does not validate for slice")
	}
}

func TestLteValForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []string `validate:"lte=2"`
	}{
		field: []string{"a", "b", "c"},
	}) {
		t.Errorf("gte validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:"lte=2"`
	}{
		field: []string{"a", "b"},
	}) {
		t.Errorf("gte validator does not validate for slice")
	}
}

func TestGteValForArray(t *testing.T) {
	if nil == Validate(struct {
		field [1]string `validate:"gte=2"`
	}{
		field: [1]string{"a"},
	}) {
		t.Errorf("gte validator does not validate for string")
	}

	if nil != Validate(struct {
		field [2]string `validate:"gte=2"`
	}{
		field: [2]string{"a", "b"},
	}) {
		t.Errorf("gte validator does not validate for string")
	}
}

func TestLteValForArray(t *testing.T) {
	if nil == Validate(struct {
		field [3]string `validate:"lte=2"`
	}{
		field: [3]string{"a", "b", "c"},
	}) {
		t.Errorf("gte validator does not validate for array")
	}

	if nil != Validate(struct {
		field [2]string `validate:"lte=2"`
	}{
		field: [2]string{"a", "b"},
	}) {
		t.Errorf("gte validator does not validate for array")
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
		t.Errorf("empty validator does not validate for slice")
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

func TestEmptyValForArray(t *testing.T) {
	if nil == Validate(struct {
		field [1]string `validate:"empty=true"`
	}{
		field: [1]string{
			"a",
		},
	}) {
		t.Errorf("empty validator does not validate for array")
	}

	if nil != Validate(struct {
		field [0]string `validate:"empty=true"`
	}{
		field: [0]string{},
	}) {
		t.Errorf("empty validator does not validate for array")
	}

	if nil == Validate(struct {
		field [0]string `validate:"empty=false"`
	}{
		field: [0]string{},
	}) {
		t.Errorf("empty validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1]string `validate:"empty=false"`
	}{
		field: [1]string{"a"},
	}) {
		t.Errorf("empty validator does not validate for array")
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

func TestFormatValForString(t *testing.T) {
	if nil == Validate(struct {
		field string `validate:"format=email"`
	}{
		field: "abc",
	}) {
		t.Errorf("format validator does not validate for string")
	}

	if nil != Validate(struct {
		field string `validate:"format=email"`
	}{
		field: "abc@example.com",
	}) {
		t.Errorf("format validator does not validate for string")
	}
}

func TestDeepValsForStruct(t *testing.T) {
	s := " "

	if nil == Validate(struct {
		field struct {
			field int `validate:"gte=0"`
		}
	}{
		field: struct {
			field int `validate:"gte=0"`
		}{field: -1},
	}) {
		t.Errorf("gte validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field int `validate:"gte=0"`
		}
	}{
		field: struct {
			field int `validate:"gte=0"`
		}{field: 0},
	}) {
		t.Errorf("gte validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field int `validate:"lte=0"`
		}
	}{
		field: struct {
			field int `validate:"lte=0"`
		}{field: 1},
	}) {
		t.Errorf("lte validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field int `validate:"lte=0"`
		}
	}{
		field: struct {
			field int `validate:"lte=0"`
		}{field: 0},
	}) {
		t.Errorf("lte validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field string `validate:"empty=true"`
		}
	}{
		field: struct {
			field string `validate:"empty=true"`
		}{field: " "},
	}) {
		t.Errorf("empty validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field string `validate:"empty=true"`
		}
	}{
		field: struct {
			field string `validate:"empty=true"`
		}{field: ""},
	}) {
		t.Errorf("empty validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field string `validate:"empty=false"`
		}
	}{
		field: struct {
			field string `validate:"empty=false"`
		}{field: ""},
	}) {
		t.Errorf("empty validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field string `validate:"empty=false"`
		}
	}{
		field: struct {
			field string `validate:"empty=false"`
		}{field: " "},
	}) {
		t.Errorf("empty validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field int `validate:"lte=0"`
		}
	}{
		field: struct {
			field int `validate:"lte=0"`
		}{field: 0},
	}) {
		t.Errorf("lte validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field *string `validate:"nil=true"`
		}
	}{
		field: struct {
			field *string `validate:"nil=true"`
		}{field: &s},
	}) {
		t.Errorf("nil validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field *string `validate:"nil=true"`
		}
	}{
		field: struct {
			field *string `validate:"nil=true"`
		}{field: nil},
	}) {
		t.Errorf("nil validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field *string `validate:"nil=false"`
		}
	}{
		field: struct {
			field *string `validate:"nil=false"`
		}{field: nil},
	}) {
		t.Errorf("nil validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field *string `validate:"nil=false"`
		}
	}{
		field: struct {
			field *string `validate:"nil=false"`
		}{field: &s},
	}) {
		t.Errorf("nil validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field int `validate:"one_of=1,2,3"`
		}
	}{
		field: struct {
			field int `validate:"one_of=1,2,3"`
		}{field: 4},
	}) {
		t.Errorf("one_of validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field int `validate:"one_of=1,2,3"`
		}
	}{
		field: struct {
			field int `validate:"one_of=1,2,3"`
		}{field: 1},
	}) {
		t.Errorf("one_of validator does not validate for struct field")
	}

	if nil == Validate(struct {
		field struct {
			field string `validate:"format=email"`
		}
	}{
		field: struct {
			field string `validate:"format=email"`
		}{field: "abc"},
	}) {
		t.Errorf("format validator does not validate for struct field")
	}

	if nil != Validate(struct {
		field struct {
			field string `validate:"format=email"`
		}
	}{
		field: struct {
			field string `validate:"format=email"`
		}{field: "abc@example.com"},
	}) {
		t.Errorf("format validator does not validate for struct field")
	}
}

func TestDeepValsForNestedStruct(t *testing.T) {
	type A struct {
		field int `validate:"gte=0"`
	}

	type B struct {
		A
		field int `validate:"gte=0"`
	}

	if nil == Validate(B{
		A: A{
			field: -1,
		},
		field: 0,
	}) {
		t.Errorf("validator does not validate for nested struct")
	}

	if nil == Validate(B{
		A: A{
			field: 0,
		},
		field: -1,
	}) {
		t.Errorf("validator does not validate for nested struct")
	}

	if nil != Validate(B{
		A: A{
			field: 0,
		},
		field: 0,
	}) {
		t.Errorf("validator does not validate for nested struct")
	}
}

func TestDeepValsForMapKeys(t *testing.T) {
	s := " "

	if nil == Validate(struct {
		field map[int]int `validate:"[gte=0]"`
	}{
		field: map[int]int{0: 0, -1: 0},
	}) {
		t.Errorf("[gte] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"[gte=0]"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf("[gte] validator does not validate map key")
	}

	if nil == Validate(struct {
		field map[int]int `validate:"[lte=0]"`
	}{
		field: map[int]int{0: 0, 1: 0},
	}) {
		t.Errorf("[lte] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"[lte=0]"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf("[lte] validator does not validate for map key")
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

	if nil == Validate(struct {
		field map[string]int `validate:"[format=email]"`
	}{
		field: map[string]int{
			"abc": 0,
		},
	}) {
		t.Errorf("[format] validator does not validate for map key")
	}

	if nil != Validate(struct {
		field map[string]int `validate:"[format=email]"`
	}{
		field: map[string]int{
			"abc@example.com": 0,
		},
	}) {
		t.Errorf("[format] validator does not validate for map key")
	}
}

func TestDeepValsForMapValues(t *testing.T) {
	s := " "

	if nil == Validate(struct {
		field map[int]int `validate:"> gte=0"`
	}{
		field: map[int]int{0: 0, 1: -1},
	}) {
		t.Errorf(">gte validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"> gte=0"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf(">gte validator does not validate map values")
	}

	if nil == Validate(struct {
		field map[int]int `validate:"> lte=0"`
	}{
		field: map[int]int{0: 0, -1: 1},
	}) {
		t.Errorf(">lte validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]int `validate:"> lte=0"`
	}{
		field: map[int]int{0: 0},
	}) {
		t.Errorf(">lte validator does not validate for map values")
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
		field map[int]*string `validate:"> [nil=false]"`
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

	if nil == Validate(struct {
		field map[int]string `validate:"> format=email"`
	}{
		field: map[int]string{
			0: "abc",
		},
	}) {
		t.Errorf(">format validator does not validate for map values")
	}

	if nil != Validate(struct {
		field map[int]string `validate:"> format=email"`
	}{
		field: map[int]string{
			0: "abc@example.com",
		},
	}) {
		t.Errorf(">format validator does not validate for map values")
	}
}

func TestDeepValsForSlice(t *testing.T) {
	if nil == Validate(struct {
		field []int `validate:">gte=0"`
	}{
		field: []int{0, -1},
	}) {
		t.Errorf(">gte validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []int `validate:">gte=0"`
	}{
		field: []int{0, 0},
	}) {
		t.Errorf(">gte validator does not validate for slice")
	}

	if nil == Validate(struct {
		field []int `validate:">lte=0"`
	}{
		field: []int{0, 1},
	}) {
		t.Errorf(">lte validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []int `validate:">lte=0"`
	}{
		field: []int{0, 0},
	}) {
		t.Errorf(">lte validator does not validate for slice")
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

	if nil == Validate(struct {
		field []string `validate:">format=email"`
	}{
		field: []string{"abc"},
	}) {
		t.Errorf(">foormat validator does not validate for slice")
	}

	if nil != Validate(struct {
		field []string `validate:">format=email"`
	}{
		field: []string{"abc@example.com"},
	}) {
		t.Errorf(">foormat validator does not validate for slice")
	}
}

func TestDeepValsForArray(t *testing.T) {
	if nil == Validate(struct {
		field [2]int `validate:">gte=0"`
	}{
		field: [2]int{0, -1},
	}) {
		t.Errorf(">gte validator does not validate for array")
	}

	if nil != Validate(struct {
		field [2]int `validate:">gte=0"`
	}{
		field: [2]int{0, 0},
	}) {
		t.Errorf(">gte validator does not validate for array")
	}

	if nil == Validate(struct {
		field [2]int `validate:">lte=0"`
	}{
		field: [2]int{0, 1},
	}) {
		t.Errorf(">lte validator does not validate for array")
	}

	if nil != Validate(struct {
		field [2]int `validate:">lte=0"`
	}{
		field: [2]int{0, 0},
	}) {
		t.Errorf(">lte validator does not validate for array")
	}

	if nil == Validate(struct {
		field [1][1]int `validate:">empty=true"`
	}{
		field: [1][1]int{
			[1]int{0},
		},
	}) {
		t.Errorf(">empty validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1][0]int `validate:">empty=true"`
	}{
		field: [1][0]int{
			[0]int{},
		},
	}) {
		t.Errorf(">empty validator does not validate for array")
	}

	if nil == Validate(struct {
		field [1][0]int `validate:">empty=false"`
	}{
		field: [1][0]int{
			[0]int{},
		},
	}) {
		t.Errorf(">empty validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1][1]int `validate:">empty=false"`
	}{
		field: [1][1]int{
			[1]int{0},
		},
	}) {
		t.Errorf(">empty validator does not validate for array")
	}

	if nil == Validate(struct {
		field [1]*int `validate:">nil=true"`
	}{
		field: [1]*int{
			new(int),
		},
	}) {
		t.Errorf(">nil validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1]*int `validate:">nil=true"`
	}{
		field: [1]*int{nil},
	}) {
		t.Errorf(">nil validator does not validate for array")
	}

	if nil == Validate(struct {
		field [1]*int `validate:">nil=false"`
	}{
		field: [1]*int{nil},
	}) {
		t.Errorf(">nil validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1]*int `validate:">nil=false"`
	}{
		field: [1]*int{new(int)},
	}) {
		t.Errorf(">nil validator does not validate for array")
	}

	if nil == Validate(struct {
		field [1]int `validate:">one_of=1,2,3"`
	}{
		field: [1]int{4},
	}) {
		t.Errorf(">one_of validator does not validate for array")
	}

	if nil != Validate(struct {
		field [3]int `validate:">one_of=1,2,3"`
	}{
		field: [3]int{1, 2, 3},
	}) {
		t.Errorf(">one_of validator does not validate for array")
	}

	if nil == Validate(struct {
		field [1]string `validate:">format=email"`
	}{
		field: [1]string{"abc"},
	}) {
		t.Errorf(">foormat validator does not validate for array")
	}

	if nil != Validate(struct {
		field [1]string `validate:">format=email"`
	}{
		field: [1]string{"abc@example.com"},
	}) {
		t.Errorf(">foormat validator does not validate for array")
	}
}

func TestDeepValsForPtr(t *testing.T) {
	gteusOne := -1
	zero := 0
	one := 1
	four := 4
	empty := ""
	notEmpty := "a"
	abc := "abc"
	email := "abc@example.com"
	onePtr := &one
	var nilPtr *int

	if nil == Validate(struct {
		field *int `validate:">gte=0"`
	}{
		field: &gteusOne,
	}) {
		t.Errorf(">gte validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:">gte=0"`
	}{
		field: &zero,
	}) {
		t.Errorf(">gte validator does not validate for pointer")
	}

	if nil == Validate(struct {
		field *int `validate:">lte=0"`
	}{
		field: &one,
	}) {
		t.Errorf(">lte validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *int `validate:">lte=0"`
	}{
		field: &zero,
	}) {
		t.Errorf(">lte validator does not validate for pointer")
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

	if nil == Validate(struct {
		field *string `validate:">format=email"`
	}{
		field: &abc,
	}) {
		t.Errorf(">format validator does not validate for pointer")
	}

	if nil != Validate(struct {
		field *string `validate:">format=email"`
	}{
		field: &email,
	}) {
		t.Errorf(">format validator does not validate for pointer")
	}
}

func TestDeepDeepVal(t *testing.T) {
	str := " "
	emptyStr := ""
	zero := 0
	gteusOne := -1

	if nil != Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{
			map[*string]*int{
				&emptyStr: &zero,
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{
			map[*string]*int{
				&emptyStr: &gteusOne,
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{
			map[*string]*int{
				&emptyStr: nil,
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{
			map[*string]*int{
				&str: &zero,
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{
			map[*string]*int{
				nil: &zero,
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{
			map[*string]*int{},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*int `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false > gte=0"`
	}{
		field: []map[*string]*int{},
	}) {
		t.Errorf("complex validator does not validate")
	}
}

func TestDeepDeepStructVal(t *testing.T) {
	emptyStr := ""

	type SubStr struct {
		field int `validate:"gte=0"`
	}

	if nil != Validate(struct {
		field []map[*string]*SubStr `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false"`
	}{
		field: []map[*string]*SubStr{
			map[*string]*SubStr{
				&emptyStr: &SubStr{field: 0},
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*SubStr `validate:"empty=false > empty=false [nil=false > empty=true] > nil=false"`
	}{
		field: []map[*string]*SubStr{
			map[*string]*SubStr{
				&emptyStr: &SubStr{field: -1},
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil != Validate(struct {
		field []map[*string]*SubStr
	}{
		field: []map[*string]*SubStr{
			map[*string]*SubStr{
				&emptyStr: &SubStr{field: 0},
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}

	if nil == Validate(struct {
		field []map[*string]*SubStr
	}{
		field: []map[*string]*SubStr{
			map[*string]*SubStr{
				&emptyStr: &SubStr{field: -1},
			},
		},
	}) {
		t.Errorf("complex validator does not validate")
	}
}

func TestIncorrectDeepVal(t *testing.T) {
	one := 1

	if nil != Validate(struct {
		field []int `validate:"gte=0"`
	}{
		field: []int{-1},
	}) {
		t.Errorf("gte validator validates one level deep for slice")
	}

	if nil != Validate(struct {
		field *int `validate:"lte=0"`
	}{
		field: &one,
	}) {
		t.Errorf("lte validator validates one level deep for pointer")
	}
}
