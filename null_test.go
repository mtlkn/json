package json

import (
	"testing"
)

func TestNull(t *testing.T) {
	s := "null"
	p := newParser([]byte(s))
	v, _, err := p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != NullType {
		t.Fail()
		t.Error("failed to parse null")
	}

	s = "nil"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad null value")
	}

	s = "nile"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad null value")
	}

	s = "{ \"f\": null }"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != ObjectType {
		t.Fail()
		t.Error("failed to parse null")
	}

	null := new(nullValue)
	if null.Value() != nil || !null.IsEmpty() {
		t.Fail()
		t.Error("null failed to implement Value interface")
	}

}
