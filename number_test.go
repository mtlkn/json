package json

import (
	"testing"
)

func TestInt(t *testing.T) {
	s := "314"
	p := newParser([]byte(s))
	v, _, err := p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != IntType || v.Value() != 314 {
		t.Fail()
		t.Error("failed to parse int")
	}

	s = "002"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != IntType || v.Value() != 2 {
		t.Fail()
		t.Error("failed to parse int")
	}

	s = "{\"f\":314x}"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad int value")
	}

	f := Int(314)
	if f.Value() != 314 || f.IsEmpty() || f.String() != "314" {
		t.Fail()
		t.Error("failed to construct int value")
	}
}

func TestFloat(t *testing.T) {
	s := "3.14"
	p := newParser([]byte(s))
	v, _, err := p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != FloatType || v.Value() != 3.14 {
		t.Fail()
		t.Error("failed to parse float")
	}

	s = ".23"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != FloatType || v.Value() != 0.23 {
		t.Fail()
		t.Error("failed to parse float")
	}

	s = "3.14e+04"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != FloatType {
		t.Fail()
		t.Error("failed to parse float")
	}

	s = "3.14E-04"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != FloatType {
		t.Fail()
		t.Error("failed to parse float")
	}

	s = "+.02"
	p = newParser([]byte(s))
	v, _, err = p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != FloatType || v.Value() != 0.02 {
		t.Fail()
		t.Error("failed to parse float")
	}

	s = "+"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad float value")
	}

	s = "3.14.14-2+3"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad float value")
	}

	s = "{\"f\":3.14x}"
	p = newParser([]byte(s))
	if _, _, err := p.ParseValue(false); err == nil {
		t.Fail()
		t.Error("failed to fail bad float value")
	}

	f := Float(3.14)
	if f.Value() != 3.14 || f.IsEmpty() || f.String() != "3.14" {
		t.Fail()
		t.Error("failed to construct float value")
	}
}

func TestUInt(t *testing.T) {
	s := "9223372036854775809"
	p := newParser([]byte(s))
	v, _, err := p.ParseValue(false)
	if (err != nil && err != errEOF) || v.Type() != UIntType {
		t.Fail()
		t.Error("failed to parse uint")
	}

	f := UInt(9223372036854775809)
	if f.Value().(uint64) != 9223372036854775809 || f.IsEmpty() || f.String() != "9223372036854775809" {
		t.Fail()
		t.Error("failed to construct uint value")
	}
}
