package cookies

import (
	"net/http"
	"testing"
)

func TestMarshal(t *testing.T) {
	type testType struct {
		FirstValue  string `cookie:"first_value"`
		SecondValue string `cookie:"second_value"`
		IntValue    int    `cookie:"int_value"`
	}

	testValue := testType{
		FirstValue:  "Hello world",
		SecondValue: "This is a test",
		IntValue:    10,
	}

	c := TypedCookie[testType]{
		TypedValue: testValue,
	}

	s, err := Marshal(c)
	if err != nil {
		t.Fatalf("Error trying to parse value:\n%s", err.Error())
	}

	expected := "first_value:Hello world|second_value:This is a test|int_value:10"

	if s.Value != expected {
		t.Fatalf("Assertion failed, expected:\n%s\n\nfound:\n%s", expected, s.Value)
	}
}

func TestUnmarshal(t *testing.T) {
	type TestType struct {
		FirstValue  string `cookie:"first_value"`
		SecondValue string `cookie:"second_value"`
	}

	c := http.Cookie{
		Value: "first_value:Hello world|second_value:This is a test",
	}

	var tc TypedCookie[TestType]
	err := Unmarshal(c, &tc)
	if err != nil {
		t.Fatalf("Error trying to parse value:\n%s", err.Error())
	}

	expected := TestType{
		FirstValue:  "Hello world",
		SecondValue: "This is a test",
	}

	if tc.TypedValue != expected {
		t.Fatalf("Assertion failed, expected:\n%s\n\nfound:\n%s", expected, tc.TypedValue)
	}
}
