package json

import (
	"reflect"
	"strings"
)

type Array struct {
	Values []*Value
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
	case []interface{}:
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

func (ja *Array) GetValue(pos int) (*Value, bool) {
	if ja == nil || len(ja.Values) <= pos {
		return nil, false
	}

	return ja.Values[pos], true
}

func (ja *Array) GetString(pos int) (string, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return "", false
	}
	return jv.GetString()
}

func (ja *Array) GetInt(pos int) (int, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return 0, false
	}
	return jv.GetInt()
}

func (ja *Array) GetUInt(pos int) (uint, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return 0, false
	}
	return jv.GetUInt()
}

func (ja *Array) GetFloat(pos int) (float64, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return 0, false
	}
	return jv.GetFloat()
}

func (ja *Array) GetBool(pos int) (bool, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return false, false
	}
	return jv.GetBool()
}

func (ja *Array) GetObject(pos int) (*Object, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return nil, false
	}
	return jv.GetObject()
}

func (ja *Array) GetArray(pos int) (*Array, bool) {
	jv, ok := ja.GetValue(pos)
	if !ok {
		return nil, false
	}
	return jv.GetArray()
}

func (ja *Array) GetStrings() ([]string, bool) {
	var vs []string
	for _, jv := range ja.Values {
		v, ok := jv.GetString()
		if !ok {
			return nil, false
		}
		vs = append(vs, v)
	}

	return vs, true
}

func (ja *Array) GetInts() ([]int, bool) {
	var vs []int
	for _, jv := range ja.Values {
		v, ok := jv.GetInt()
		if !ok {
			return nil, false
		}
		vs = append(vs, v)
	}

	return vs, true
}

func (ja *Array) GetFloats() ([]float64, bool) {
	var vs []float64
	for _, jv := range ja.Values {
		v, ok := jv.GetFloat()
		if !ok {
			return nil, false
		}
		vs = append(vs, v)
	}

	return vs, true
}

func (ja *Array) GetObjects() ([]*Object, bool) {
	var vs []*Object
	for _, jv := range ja.Values {
		v, ok := jv.GetObject()
		if !ok {
			return nil, false
		}
		vs = append(vs, v)
	}

	return vs, true
}

func (ja *Array) Add(value interface{}) *Array {
	ja.Values = append(ja.Values, newValue(value))
	return ja
}

func (ja *Array) Remove(pos int) *Array {
	if pos < len(ja.Values) {
		var vs []*Value
		vs = append(vs, ja.Values[:pos]...)
		vs = append(vs, ja.Values[pos+1:]...)
		ja.Values = vs
	}

	return ja
}

func (ja *Array) Merge(merge *Array) *Array {
	if merge == nil || len(merge.Values) == 0 {
		return ja
	}

	for _, jv := range merge.Values {
		v, err := jv.GetValue()
		if err != nil && v != nil {
			ja.Add(v)
		}
	}

	return ja
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

func (ja *Array) Bytes() []byte {
	return stringToBytes(ja.String())
}
