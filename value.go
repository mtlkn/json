package json

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ValueType uint8

const (
	STRING ValueType = iota + 1
	INT
	UINT
	FLOAT
	BOOL
	OBJECT
	ARRAY
	NULL
)

func (vt ValueType) String() string {
	switch vt {
	case STRING:
		return "string"
	case INT:
		return "int"
	case UINT:
		return "uint"
	case FLOAT:
		return "float"
	case BOOL:
		return "bool"
	case OBJECT:
		return "object"
	case ARRAY:
		return "array"
	case NULL:
		return "null"
	default:
		return "unknown"
	}
}

type specialBytes uint8

const (
	unquoteBytes specialBytes = iota + 1
	floatBytes
)

type Value struct {
	Type    ValueType
	value   interface{}
	data    []byte
	special specialBytes
}

func (jv *Value) String() string {
	switch jv.Type {
	case OBJECT:
		jo, ok := jv.GetObject()
		if !ok {
			return ""
		}
		return jo.String()
	case ARRAY:
		ja, ok := jv.GetArray()
		if !ok {
			return ""
		}
		return ja.String()
	default:
		if len(jv.data) == 0 {
			return jv.toString()
		}
		return bytesToString(jv.data)
	}
}

func (jv *Value) GetValue() (interface{}, error) {
	if jv.value != nil {
		return jv.value, nil
	}

	if jv.Type == STRING {
		if jv.special != unquoteBytes {
			jv.value = bytesToString(jv.data[1 : len(jv.data)-1])
		} else {
			s, err := strconv.Unquote(bytesToString(jv.data))
			if err != nil {
				return "", err
			}
			jv.value = s
		}

		return jv.value, nil
	}

	if len(jv.data) == 0 {
		return nil, errors.New("empty")
	}

	if jv.data[0] == 't' && len(jv.data) == 4 && jv.data[1] == 'r' && jv.data[2] == 'u' && jv.data[3] == 'e' {
		jv.Type = BOOL
		jv.value = true
		return jv.value, nil
	}

	if jv.data[0] == 'f' && len(jv.data) == 5 && jv.data[1] == 'a' && jv.data[2] == 'l' && jv.data[3] == 's' && jv.data[4] == 'e' {
		jv.Type = BOOL
		jv.value = false
		return jv.value, nil
	}

	if jv.data[0] == 'n' && len(jv.data) == 4 && jv.data[1] == 'u' && jv.data[2] == 'l' && jv.data[3] == 'l' {
		jv.Type = NULL
		return jv.value, nil
	}

	s := bytesToString(jv.data)

	if jv.special != floatBytes {
		d, err := strconv.Atoi(s)
		if err != nil {
			if strings.Contains(err.Error(), "value out of range") && d > 0 {
				u, err := strconv.ParseUint(s, 10, 0)
				if err == nil {
					jv.Type = UINT
					jv.value = uint(u)
					return jv.value, nil
				}
			}
			return nil, err
		}

		jv.Type = INT
		jv.value = d
		return d, nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New("non-json value")
	}

	jv.Type = FLOAT
	jv.value = f
	return f, nil
}

func (jv *Value) GetString() (string, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != STRING {
		return "", false
	}

	s, ok := v.(string)
	return s, ok
}

func (jv *Value) GetInt() (int, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != INT {
		return 0, false
	}

	i, ok := v.(int)
	return i, ok
}

func (jv *Value) GetUInt() (uint, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != UINT {
		return 0, false
	}

	ui, ok := v.(uint)
	return ui, ok
}

func (jv *Value) GetFloat() (float64, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != FLOAT {
		return 0, false
	}

	f, ok := v.(float64)
	return f, ok
}

func (jv *Value) GetBool() (bool, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != BOOL {
		return false, false
	}

	b, ok := v.(bool)
	return b, ok
}

func (jv *Value) GetObject() (*Object, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != OBJECT {
		return nil, false
	}

	jo, ok := v.(*Object)
	return jo, ok
}

func (jv *Value) GetArray() (*Array, bool) {
	v, err := jv.GetValue()
	if err != nil || v == nil || jv.Type != ARRAY {
		return nil, false
	}

	ja, ok := v.(*Array)
	return ja, ok
}

func (jv *Value) toString() string {
	if jv.value == nil {
		return "null"
	}

	switch jv.Type {
	case STRING:
		s := strconv.Quote(jv.value.(string))
		jv.data = stringToBytes(s)
		return s
	case INT:
		s := strconv.Itoa(jv.value.(int))
		jv.data = stringToBytes(s)
		return s
	case UINT:
		s := strconv.FormatUint(uint64(jv.value.(uint)), 10)
		jv.data = stringToBytes(s)
		return s
	case FLOAT:
		s := strconv.FormatFloat(jv.value.(float64), 'f', -1, 64)
		jv.data = stringToBytes(s)
		return s
	case BOOL:
		v := jv.value.(bool)
		if v {
			return "true"
		}
		return "false"
	}

	return "null"
}

func newValue(x interface{}) *Value {
	var t ValueType
	v := x

	switch x := x.(type) {
	case string:
		t = STRING
	case int:
		t = INT
	case uint:
		t = UINT
	case float64:
		t = FLOAT
	case bool:
		t = BOOL
	case *Object:
		t = OBJECT
	case *Array:
		t = ARRAY
	case int8:
		t = INT
		v = int(x)
	case int16:
		t = INT
		v = int(x)
	case int32:
		t = INT
		v = int(x)
	case int64:
		t = INT
		v = int(x)
	case uint8:
		t = UINT
		v = uint(x)
	case uint16:
		t = UINT
		v = uint(x)
	case uint32:
		t = UINT
		v = uint(x)
	case uint64:
		t = UINT
		v = uint(x)
	case float32:
		t = FLOAT
		v, _ = strconv.ParseFloat(fmt.Sprintf("%v", x), 64)
	case []string:
		t = ARRAY
		v = NewArray(x)
	case []int, []int8, []int16, []int32, []int64:
		t = ARRAY
		v = NewArray(x)
	case []uint, []uint8, []uint16, []uint32, []uint64:
		t = ARRAY
		v = NewArray(x)
	case []float64, []float32:
		t = ARRAY
		v = NewArray(x)
	case []bool:
		t = ARRAY
		v = NewArray(x)
	case []*Object, []*Array:
		t = ARRAY
		v = NewArray(x)
	}

	if t == 0 {
		return &Value{
			Type: NULL,
		}
	}

	return &Value{
		Type:  t,
		value: v,
	}
}
