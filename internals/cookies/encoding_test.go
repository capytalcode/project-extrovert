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
}

var TEST_VALUE = TestType{
	FirstValue:   "Hello world",
	SecondValue:  "This is a test",
	IntValue:     -10,
}

var TEST_TYPED_COOKIE = TypedCookie[TestType]{
	TypedValue: TEST_VALUE,
}

var TEST_STRING = strings.Join([]string{
	"first_value:Hello world",
	"second_value:This is a test",
	"int_value:-10",
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
