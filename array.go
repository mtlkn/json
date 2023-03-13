package json

import (
	"reflect"
	"strings"
)

type Array struct {
	Values []*Value
}

func (ja *Array) String() string {
	if ja == nil {
		return ""
	}

	if len(ja.Values) == 0 {
		return "[]"
	}

	var sb strings.Builder

	sb.WriteByte('[')

	for i, jv := range ja.Values {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(jv.String())
	}

	sb.WriteByte(']')

	return sb.String()
}

func NewArray(x interface{}) *Array {
	var vs []*Value

	switch x := x.(type) {
	case []string:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []int:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []uint:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []float64:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []bool:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []*Object:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []*Array:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []int8:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []int16:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []int32:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []int64:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []uint8:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []uint16:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []uint32:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []uint64:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	case []float32:
		for _, v := range x {
			vs = append(vs, newValue(v))
		}
	default:
		if x == nil {
			return nil
		}
		if reflect.TypeOf(x).Kind() != reflect.Slice {
			v := newValue(x)
			if v.Type > 0 {
				vs = []*Value{v}
			}
		}
	}

	if len(vs) == 0 {
		return nil
	}

	return &Array{
		Values: vs,
	}
}
