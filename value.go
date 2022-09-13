package json

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ValueType int

const (
	OBJECT ValueType = iota + 1
	ARRAY
	STRING
	INT
	UINT
	FLOAT
	BOOL
	NULL
	INVALID
)

type Value struct {
	typ ValueType
	buf []byte      // parsed bytes
	val interface{} // actual value
}

func New(v interface{}) *Value {
	switch v.(type) {
	case string:
		return &Value{
			typ: STRING,
			val: v,
		}
	case int:
		return &Value{
			typ: INT,
			val: v,
		}
	case uint:
		return &Value{
			typ: UINT,
			val: v,
		}
	case float64:
		return &Value{
			typ: FLOAT,
			val: v,
		}
	case bool:
		return &Value{
			typ: BOOL,
			val: v,
		}
	case *Object:
		return &Value{
			typ: OBJECT,
			val: v,
		}
	case *Array:
		return &Value{
			typ: ARRAY,
			val: v,
		}
	case int8:
		return &Value{
			typ: INT,
			val: int((v).(int8)),
		}
	case int16:
		return &Value{
			typ: INT,
			val: int((v).(int16)),
		}
	case int32:
		return &Value{
			typ: INT,
			val: int((v).(int32)),
		}
	case int64:
		return &Value{
			typ: INT,
			val: int((v).(int64)),
		}
	case uint8:
		return &Value{
			typ: UINT,
			val: uint((v).(uint8)),
		}
	case uint16:
		return &Value{
			typ: UINT,
			val: uint((v).(uint16)),
		}
	case uint32:
		return &Value{
			typ: UINT,
			val: uint((v).(uint32)),
		}
	case uint64:
		return &Value{
			typ: UINT,
			val: uint((v).(uint64)),
		}
	case float32:
		return &Value{
			typ: FLOAT,
			val: Float32To64((v).(float32)),
		}
	}

	return nil
}

func (v *Value) Debug() string {
	return fmt.Sprintf("{\"type\":%v,\"parsed\":\"%s\",\"value\":%v}", v.typ, BytesToString(v.buf), v.val)
}

func (v *Value) Type() ValueType {
	if v.typ == 0 {
		v.Validate()
	}
	return v.typ
}

func (v *Value) Value() interface{} {
	if v.typ == 0 {
		v.Validate()
	}
	return v.val
}

func (v *Value) String() (string, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	if v.typ != STRING {
		return "", false
	}

	s, ok := (v.val).(string)
	return s, ok
}

func (v *Value) Int() (int, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	switch v.typ {
	case INT:
		i, ok := (v.val).(int)
		return i, ok
	case FLOAT:
		f, ok := (v.val).(float64)
		return int(f), ok
	}

	return 0, false
}

func (v *Value) UInt() (uint, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	switch v.typ {
	case UINT:
		u, ok := (v.val).(uint)
		return u, ok
	case INT:
		i, ok := (v.val).(int)
		return uint(i), ok
	case FLOAT:
		f, ok := (v.val).(float64)
		return uint(f), ok
	}

	return 0, false
}

func (v *Value) Float() (float64, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	switch v.typ {
	case FLOAT:
		f, ok := (v.val).(float64)
		return f, ok
	case INT:
		i, ok := (v.val).(int)
		return float64(i), ok
	}

	return 0, false
}

func (v *Value) Bool() (bool, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	if v.typ != BOOL {
		return false, false
	}

	b, ok := (v.val).(bool)
	return b, ok
}

func (v *Value) Object() (*Object, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	if v.typ != OBJECT {
		return nil, false
	}

	jo, ok := (v.val).(*Object)
	return jo, ok
}

func (v *Value) Array() (*Array, bool) {
	if v.typ == 0 {
		v.Validate()
	}

	if v.typ != ARRAY {
		return nil, false
	}

	ja, ok := (v.val).(*Array)
	return ja, ok
}

// we need to validate parsed JSON only, constucted JSON is always valid
func (v *Value) Validate() error {
	if v == nil {
		return errors.New("nil value")
	}

	if v.typ > 0 && v.typ < INVALID && v.val != nil {
		return nil
	}

	if len(v.buf) == 0 {
		v.typ = INVALID
		return errors.New("missing value")
	}

	switch v.buf[0] {
	case '"':
		if v.parseString() {
			return nil
		}
	case 't':
		if len(v.buf) == 4 && v.buf[1] == 'r' && v.buf[2] == 'u' && v.buf[3] == 'e' {
			v.typ = BOOL
			v.val = true
			return nil
		}
	case 'f':
		if len(v.buf) == 5 && v.buf[1] == 'a' && v.buf[2] == 'l' && v.buf[3] == 's' && v.buf[4] == 'e' {
			v.typ = BOOL
			v.val = false
			return nil
		}
	case 'n':
		if len(v.buf) == 4 && v.buf[1] == 'u' && v.buf[2] == 'l' && v.buf[3] == 'l' {
			v.typ = NULL
			return nil
		}
	default:
		if v.parseNumber() {
			return nil
		}
	}

	v.typ = INVALID
	return errors.New("invalid JSON value: " + BytesToString(v.buf))
}

func (v *Value) parseString() bool {
	s := BytesToString(v.buf)
	s, err := strconv.Unquote(s)
	if err != nil {
		v.typ = INVALID
		return false
	}

	v.val = s
	v.typ = STRING
	return true
}

func (v *Value) parseNumber() bool {
	var float bool

	for i, b := range v.buf {
		if (b >= '0' && b <= '9') || ((b == '-' || b == '+') && i == 0) {
			continue
		} else if b == '.' || b == 'e' || b == 'E' || b == '-' || b == '+' {
			float = true
		} else {
			v.typ = INVALID
			return false
		}
	}

	s := BytesToString(v.buf)

	if float {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			v.val = f
			v.typ = FLOAT
			return true
		}
		return false
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		if strings.Contains(err.Error(), "value out of range") && i > 0 {
			u, err := strconv.ParseUint(s, 10, 0)
			if err == nil {
				v.val = uint(u)
				v.typ = UINT
				return true
			}
		}

		v.typ = INVALID
		return false
	}

	v.val = i
	v.typ = INT
	return true
}

func (v *Value) string() string {
	if v == nil {
		return "null"
	}

	switch v.Type() {
	case STRING:
		s, _ := v.String()
		return strconv.Quote(s)
	case OBJECT:
		o, _ := v.Object()
		if o == nil {
			return "null"
		}
		return o.String()
	case ARRAY:
		a, _ := v.Array()
		if a == nil {
			return "null"
		}
		return a.String()
	case INT:
		i, _ := v.Int()
		return strconv.Itoa(i)
	case FLOAT:
		f, _ := v.Float()
		return strconv.FormatFloat(f, 'f', -1, 64)
	case BOOL:
		if v.Value() == true {
			return "true"
		}
		return "false"
	case UINT:
		u, _ := v.UInt()
		return strconv.FormatUint(uint64(u), 10)
	}

	return "null"
}
