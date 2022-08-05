package json

import (
	"testing"
)

func TestBool(t *testing.T) {
	s := "true"
	p := newParser([]byte(s))
	v, _, err := p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != BoolType {
		t.Fail()
		t.Error("failed to parse bool true")
	}

	s = "tru"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad bool true value")
	}

	s = "tree"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad bool true value")
	}

	s = "false"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != BoolType {
		t.Fail()
		t.Error("failed to parse bool false")
	}

	s = "fals"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad bool false value")
	}

	s = "faust"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad bool false value")
	}

	s = "{ \"f\": true }"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != ObjectType {
		t.Fail()
		t.Error("failed to parse bool true")
	}

	b := Bool(true)
	if b.Value() != true || b.IsEmpty() || b.String() != "true" {
		t.Fail()
		t.Error("failed to construct bool value")
	}

	b = Bool(false)
	if b.Value() != false || !b.IsEmpty() || b.String() != "false" {
		t.Fail()
		t.Error("failed to construct bool value")
	}
}
