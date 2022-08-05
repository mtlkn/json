package json

import (
	"testing"
)

func TestArray(t *testing.T) {
	ja := NewArray(Int(1), Int(2))
	if len(ja.Values) != 2 && ja.Values[1].Value() != 2 {
		t.Fail()
		t.Error("failed to init new array")
	}

	ja.AddInt(3)
	if len(ja.Values) != 3 && ja.Values[2].Value() != 3 {
		t.Fail()
		t.Error("failed to add int to array")
	}

	ja.AddFloat(3.14)
	if len(ja.Values) != 4 && ja.Values[3].Value() != 3.14 {
		t.Fail()
		t.Error("failed to add float to array")
	}

	ja.AddString("four")
	if len(ja.Values) != 5 && ja.Values[4].Value() != "four" {
		t.Fail()
		t.Error("failed to add string to array")
	}

	ja.AddObject(New())
	if len(ja.Values) != 6 {
		t.Fail()
		t.Error("failed to add object to array")
	}

	if _, ok := ja.GetStrings(); ok {
		t.Fail()
		t.Error("failed to fail to get all strings from array")
	}

	if _, ok := ja.GetInts(); ok {
		t.Fail()
		t.Error("failed to fail to get all ints from array")
	}

	if _, ok := ja.GetFloats(); ok {
		t.Fail()
		t.Error("failed to fail to get all floats from array")
	}

	if _, ok := ja.GetObjects(); ok {
		t.Fail()
		t.Error("failed to fail to get all objects from array")
	}

	ja = nil
	if ja.String() != "[]" {
		t.Fail()
		t.Error("failed to string nil array")
	}
	if ja.Copy() != nil {
		t.Fail()
		t.Error("failed to copy nil array")
	}

	ja = NewArray()
	if ja.String() != "[]" {
		t.Fail()
		t.Error("failed to string empty array")
	}
	ja = ja.Copy()
	if len(ja.Values) > 0 {
		t.Fail()
		t.Error("failed to fail empty array")
	}

	s := `[ "one", "${id?two}", 3, true ]`
	ja, _ = ParseArrayWithParameters([]byte(s))
	ja = ja.Copy()
	if len(ja.Values) != 4 || ja.Values[1].Value() != "${id?two}" {
		t.Fail()
		t.Error("failed to fail parameterized array")
	}

	s = `[ "${one}", { "f": "${two}" }, [ "${three}" ] ]`
	ja, _ = ParseArrayWithParameters([]byte(s))
	ps := ja.GetParameters()
	if len(ps) != 3 {
		t.Fail()
		t.Error("failed to get array parameters")
	}

	var l, r *Array
	if ok, _ := l.Equals(r); !ok {
		t.Fail()
		t.Error("failed to compare two nil arrays")
	}

	l = NewArray()
	r = NewArray(Int(1))
	if ok, _ := l.Equals(r); ok {
		t.Fail()
		t.Error("failed to fail compare two mismatched arrays")
	}

	ja = NewArray(Int(1))
	if ja.Type() != ArrayType {
		t.Fail()
		t.Error("failed to have propert array type")
	}
	v := ja.Value().([]Value)
	if v[0].Value() != 1 {
		t.Fail()
		t.Error("failed to compare array value to itself")
	}
}

func TestArrayParse(t *testing.T) {
	s := "[ 1, 2, 3 ]"
	ja, err := ParseArray([]byte(s))
	if err != nil || len(ja.Values) != 3 || ja.Values[1].Value() != 2 {
		t.Error("Failed to parse array")
	}

	s = `[ "one", "${id?two}", 3, true ]`
	ja, err = ParseArrayWithParameters([]byte(s))
	if err != nil || len(ja.Values) != 4 {
		t.Error("Failed to parse array")
	}
	ja = ja.SetParameters(New())
	if ja.Values[1].Value() != "two" {
		t.Error("Failed to set array parameter")
	}

	s = "a"
	if _, err := parseArray([]byte(s), false, true); err == nil {
		t.Error("Failed to faile to parse bar array")
	}

	s = "{}"
	if _, err := parseArray([]byte(s), false, false); err == nil {
		t.Error("Failed to faile to parse bar array")
	}

	s = "[ { 2 ]"
	if _, err := parseArray([]byte(s), false, false); err == nil {
		t.Error("Failed to faile to parse bar array")
	}
}
