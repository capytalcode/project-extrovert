package cookies

import (
	"net/http"
	"strings"
	"testing"
)

type TestType struct {
	FirstValue   string    `cookie:"first_value"`
	SecondValue  string    `cookie:"second_value"`
	IntValue     int       `cookie:"int_value"`
	UintValue    uint      `cookie:"uint_value"`
	BoolValue    bool      `cookie:"bool_value"`
	FloatValue   float64   `cookie:"float_value"`
	ComplexValue complex64 `cookie:"complex_value"`
}

var TEST_VALUE = TestType{
	FirstValue:   "Hello world",
	SecondValue:  "This is a test",
	IntValue:     -10,
	UintValue:    13,
	BoolValue:    true,
	FloatValue:   3.14159,
	ComplexValue: -43110.70519,
}

var TEST_TYPED_COOKIE = TypedCookie[TestType]{
	TypedValue: TEST_VALUE,
}

var TEST_STRING = strings.Join([]string{
	"first_value:Hello world",
	"second_value:This is a test",
	"int_value:-10",
	"uint_value:13",
	"bool_value:true",
	"float_value:3.14159",
	"complex_value:(-43110.707+0i)",
}, "|")

func TestMarshal(t *testing.T) {
	s, err := Marshal(TEST_TYPED_COOKIE)
	if err != nil {
		t.Fatalf("Error trying to parse value:\n%s", err.Error())
	}

	if s.Value != TEST_STRING {
		t.Fatalf("Assertion failed, expected:\n%s\n\nfound:\n%s", TEST_STRING, s.Value)
	}
}

func TestUnmarshal(t *testing.T) {
	var tc TypedCookie[TestType]
	err := Unmarshal(http.Cookie{Value: TEST_STRING}, &tc)
	if err != nil {
		t.Fatalf("Error trying to parse value:\n%s", err.Error())
	}

	if tc.TypedValue != TEST_VALUE {
		t.Fatalf("Assertion failed, expected:\n%v\n\nfound:\n%v", TEST_VALUE, tc.TypedValue)
	}
}
