package json

import (
	"testing"
)

func TestStringParse(t *testing.T) {
	s := "\"joe biden\""
	p := newParser([]byte(s))
	v, _, err := p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != StringType || v.Value() != "joe biden" {
		t.Fail()
		t.Error("failed to parse string")
	}

	s = "\"joe \\n \\r \\t \\a \\b \\f \\v \\\\\\\"xyz\\\\\\\" \\ym \\\"biden\\\"\""
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != StringType {
		t.Fail()
		t.Error("failed to parse string")
	}

	s = "\"\\u0059\\u004D\""
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != StringType {
		t.Fail()
		t.Error("failed to parse string")
	}

	s = "\"joe biden"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad string value")
	}

	s = "\"\\uxv0059\\u004D\""
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad string value")
	}

	s = "\"_search/${index\""
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(true); err == nil {
		t.Fail()
		t.Error("failed to fail bad parameterized string value")
	}
}
