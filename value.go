package json

import (
	"errors"
	"fmt"
)

type ValueType int

// Value types
const (
	ObjectType ValueType = iota + 1
	ArrayType
	StringType
	IntType
	UIntType
	FloatType
	BoolType
	NullType
)

func (vt ValueType) String() string {
	switch vt {
	case ObjectType:
		return "object"
	case ArrayType:
		return "array"
	case StringType:
		return "string"
	case IntType:
		return "int"
	case UIntType:
		return "uint64"
	case FloatType:
		return "float64"
	case BoolType:
		return "bool"
	}

	return "null"
}

// JSON value interface
type Value interface {
	Value() interface{}
	Type() ValueType
	String() string
	IsEmpty() bool
}

func (p *byteParser) ParseValue(parameterized bool) (Value, []Parameter, error) {
	err := p.SkipWS()
	if err != nil {
		return nil, nil, err
	}

	switch p.Byte {
	case '"':
		return p.ParseString(parameterized)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+', '.':
		v, err := p.ParseNumber()
		return v, nil, err
	case '{':
		jo, err := p.ParseObject(parameterized)
		return jo, nil, err
	case '[':
		ja, err := p.ParseArray(parameterized)
		return ja, nil, err
	case 't':
		v, err := p.ParseTrue()
		return v, nil, err
	case 'f':
		v, err := p.ParseFalse()
		return v, nil, err
	case 'n':
		v, err := p.ParseNull()
		return v, nil, err
	}

	return nil, nil, errors.New("invalid JSON")
}

func ObjectValue(v Value) (*Object, bool) {
	if v.Type() != ObjectType {
		return nil, false
	}
	return v.(*Object), true
}

func ArrayValue(v Value) (*Array, bool) {
	if v.Type() != ArrayType {
		return nil, false
	}
	return v.(*Array), true
}

func StringValue(v Value) (string, bool) {
	if v.Type() != StringType {
		return "", false
	}
	return (v.Value()).(string), true
}

func IntValue(v Value) (int, bool) {
	switch v.Type() {
	case IntType:
		return (v.Value()).(int), true
	case FloatType:
		f := (v.Value()).(float64)
		return int(f), true
	}

	return 0, false
}

func FloatValue(v Value) (float64, bool) {
	switch v.Type() {
	case FloatType:
		return (v.Value()).(float64), true
	case IntType:
		i := (v.Value()).(int)
		return float64(i), true
	}

	return 0, false
}

func BoolValue(v Value) (bool, bool) {
	if v.Type() != BoolType {
		return false, false
	}
	return (v.Value()).(bool), true
}

func NullValue(v Value) bool {
	return v.Type() == NullType
}

func copyValue(v Value) (Value, bool) {
	switch v.Type() {
	case StringType, IntType, FloatType, BoolType:
		return v, true
	case ObjectType:
		jo, _ := ObjectValue(v)
		copy := jo.Copy()
		return copy, true
	case ArrayType:
		ja, _ := ArrayValue(v)
		copy := ja.Copy()
		return copy, true
	}

	return nil, false
}

func compareValues(left Value, right Value) error {
	lt := left.Type()
	rt := right.Type()
	if lt != rt {
		if ((lt == IntType || lt == UIntType) && rt == FloatType) || (lt == FloatType && (rt == IntType || rt == UIntType)) {
			lt = FloatType
		} else {
			return fmt.Errorf("different types: %s != %s", lt.String(), rt.String())
		}
	}

	switch lt {
	case StringType:
		l, _ := StringValue(left)
		r, _ := StringValue(right)
		if l != r {
			return fmt.Errorf("\"%s\" != \"%s\"", l, r)
		}
	case IntType, UIntType:
		l, _ := IntValue(left)
		r, _ := IntValue(right)
		if l != r {
			return fmt.Errorf("%d != %d", l, r)
		}
	case BoolType:
		l, _ := BoolValue(left)
		r, _ := BoolValue(right)
		if l != r {
			return fmt.Errorf("%v != %v", l, r)
		}
	case FloatType:
		l, _ := FloatValue(left)
		r, _ := FloatValue(right)
		if l != r {
			return fmt.Errorf("%v != %v", l, r)
		}
	case ObjectType:
		l, _ := ObjectValue(left)
		r, _ := ObjectValue(right)
		_, err := l.Equals(r)
		if err != nil {
			return err
		}
	case ArrayType:
		l, _ := ArrayValue(left)
		r, _ := ArrayValue(right)
		_, err := l.Equals(r)
		if err != nil {
			return err
		}
	}

	return nil
}
