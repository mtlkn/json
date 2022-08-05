package json

import "testing"

func TestValue(t *testing.T) {
	p := newParser([]byte{})
	if _, _, err := p.ParseValue(false); err != errEOF {
		t.Fail()
		t.Error("failed to fail to parse empty bytes")
	}

	l := Int(1)
	r := Float(1.0)
	if err := compareValues(l, r); err != nil {
		t.Fail()
		t.Error("failed to compare same int and float")
	}

	r = String("1")
	if err := compareValues(l, r); err == nil {
		t.Fail()
		t.Error("failed to faile to compare int and string")
	}

	r = Int(2)
	if err := compareValues(l, r); err == nil {
		t.Fail()
		t.Error("failed to faile to compare 1 and 2")
	}

	l = Bool(true)
	r = Bool(false)
	if err := compareValues(l, r); err == nil {
		t.Fail()
		t.Error("failed to faile to compare true and false")
	}

	l = Float(3.14)
	r = Float(0.14)
	if err := compareValues(l, r); err == nil {
		t.Fail()
		t.Error("failed to faile to compare 3.14 and 0.14")
	}

	l = NewIntArray([]int{1})
	r = NewIntArray([]int{2})
	if err := compareValues(l, r); err == nil {
		t.Fail()
		t.Error("failed to faile to compare [1]] and [2]")
	}

	l = String("ym")
	scp, ok := copyValue(l)
	if !ok || scp.Value() != "ym" {
		t.Fail()
		t.Error("failed to copy string")
	}

	l = Int(1)
	icp, ok := copyValue(l)
	if !ok || icp.Value() != 1 {
		t.Fail()
		t.Error("failed to copy int")
	}

	l = New(Field("id", Int(1)))
	if _, ok := copyValue(l); !ok {
		t.Fail()
		t.Error("failed to copy object")
	}

	l = NewIntArray([]int{1})
	if _, ok := copyValue(l); !ok {
		t.Fail()
		t.Error("failed to copy array")
	}

}

func TestValueToString(t *testing.T) {
	if ObjectType.String() != "object" {
		t.Fail()
	}
	if ArrayType.String() != "array" {
		t.Fail()
	}
	if StringType.String() != "string" {
		t.Fail()
	}
	if IntType.String() != "int" {
		t.Fail()
	}
	if FloatType.String() != "float64" {
		t.Fail()
	}
	if UIntType.String() != "uint64" {
		t.Fail()
	}
	if BoolType.String() != "bool" {
		t.Fail()
	}
	if NullType.String() != "null" {
		t.Fail()
	}
}

func TestValueConstructors(t *testing.T) {
	var (
		jo Value = new(Object)
		ja Value = new(Array)
		s  Value = new(stringValue)
		i  Value = new(intValue)
		f  Value = new(floatValue)
		b  Value = new(boolValue)
		n  Value = new(nullValue)
	)

	if _, ok := ObjectValue(jo); !ok {
		t.Fail()
	}
	if _, ok := ObjectValue(s); ok {
		t.Fail()
	}

	if _, ok := ArrayValue(ja); !ok {
		t.Fail()
	}
	if _, ok := ArrayValue(s); ok {
		t.Fail()
	}

	if _, ok := StringValue(s); !ok {
		t.Fail()
	}
	if _, ok := StringValue(i); ok {
		t.Fail()
	}

	if _, ok := IntValue(i); !ok {
		t.Fail()
	}
	if _, ok := IntValue(f); !ok {
		t.Fail()
	}
	if _, ok := IntValue(s); ok {
		t.Fail()
	}

	if _, ok := FloatValue(f); !ok {
		t.Fail()
	}
	if _, ok := FloatValue(i); !ok {
		t.Fail()
	}
	if _, ok := FloatValue(s); ok {
		t.Fail()
	}

	if _, ok := BoolValue(b); !ok {
		t.Fail()
	}
	if _, ok := BoolValue(s); ok {
		t.Fail()
	}

	if ok := NullValue(n); !ok {
		t.Fail()
	}
	if ok := NullValue(s); ok {
		t.Fail()
	}
}
