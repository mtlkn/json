package json

import "testing"

func TestProperty(t *testing.T) {
	jp := &Property{
		Value: String("ym"),
	}
	s, ok := jp.GetString()
	if !ok || s != "ym" {
		t.Fail()
		t.Error("failed to get property string value")
	}
	if _, ok := jp.GetStrings(); ok {
		t.Fail()
		t.Error("failed to fail to get property strings value")
	}
	jp.Value = NewStringArray([]string{"ym"})
	if _, ok := jp.GetStrings(); !ok {
		t.Fail()
		t.Error("failed to get property strings value")
	}

	jp.Value = Int(1)
	i, ok := jp.GetInt()
	if !ok || i != 1 {
		t.Fail()
		t.Error("failed to get property int value")
	}
	if _, ok := jp.GetInts(); ok {
		t.Fail()
		t.Error("failed to fail to get property ints value")
	}
	jp.Value = NewIntArray([]int{1})
	if _, ok := jp.GetInts(); !ok {
		t.Fail()
		t.Error("failed to get property ints value")
	}

	jp.Value = Float(3.14)
	f, ok := jp.GetFloat()
	if !ok || f != 3.14 {
		t.Fail()
		t.Error("failed to get property float value")
	}
	if _, ok := jp.GetFloats(); ok {
		t.Fail()
		t.Error("failed to fail to get property floats value")
	}
	jp.Value = NewFloatArray([]float64{1})
	if _, ok := jp.GetFloats(); !ok {
		t.Fail()
		t.Error("failed to get property floats value")
	}

	jp.Value = Bool(true)
	b, ok := jp.GetBool()
	if !ok || !b {
		t.Fail()
		t.Error("failed to get property bool value")
	}

	jp.Value = New(Field("id", Int(1)))
	jo, ok := jp.GetObject()
	if !ok {
		t.Fail()
		t.Error("failed to get property object value")
	}
	if _, ok := jp.GetObjects(); ok {
		t.Fail()
		t.Error("failed to fail to get property objects value")
	}
	jp.Value = NewObjectArray([]*Object{jo})
	if _, ok := jp.GetObjects(); !ok {
		t.Fail()
		t.Error("failed to get property objects value")
	}
}

func TestParseProperty(t *testing.T) {
	s := "{\"f\":1}"
	if _, err := ParseObject([]byte(s)); err != nil {
		t.Fail()
		t.Error("failed to parse property")
	}

	s = "{\"f}"
	if _, err := ParseObject([]byte(s)); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad property")
	}

	s = "{\"f\"}"
	if _, err := ParseObject([]byte(s)); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad property")
	}

	s = "{\"f\":}"
	if _, err := ParseObject([]byte(s)); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad property")
	}

	s = "{\"f$"
	if _, err := ParseObjectWithParameters([]byte(s)); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad property")
	}
}
