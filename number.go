package json

import (
	"fmt"
	"strconv"
	"strings"
)

// int
type intValue struct {
	value int
	text  string
}

func Int(v int) Value {
	return &intValue{
		value: v,
	}
}

func (v *intValue) Value() interface{} {
	return v.value
}

func (v *intValue) Type() ValueType {
	return IntType
}

func (v *intValue) String() string {
	if v.text == "" {
		v.text = strconv.Itoa(v.value)
	}
	return v.text
}

func (v *intValue) IsEmpty() bool {
	return v == nil || v.value == 0
}

// float
type floatValue struct {
	value float64
	text  string
}

func Float(v float64) Value {
	return &floatValue{
		value: v,
	}
}

func (v *floatValue) Value() interface{} {
	return v.value
}

func (v *floatValue) Type() ValueType {
	return FloatType
}

func (v *floatValue) String() string {
	if v.text == "" {
		v.text = fmt.Sprintf("%v", v.value)
	}
	return v.text
}

func (v *floatValue) IsEmpty() bool {
	return v == nil || v.value == 0
}

// uint
type uintValue struct {
	value uint64
	text  string
}

func UInt(v uint64) Value {
	return &uintValue{
		value: v,
	}
}

func (v *uintValue) Value() interface{} {
	return v.value
}

func (v *uintValue) Type() ValueType {
	return UIntType
}

func (v *uintValue) String() string {
	if v.text == "" {
		v.text = strconv.FormatUint(v.value, 10)
	}
	return v.text
}

func (v *uintValue) IsEmpty() bool {
	return v == nil || v.value == 0
}

//returns float64 or int
func (p *byteParser) ParseNumber() (Value, error) {
	var (
		sb    strings.Builder
		float = p.Byte == '.'
		idx   = p.Index
		err   error
	)

	sb.WriteByte(p.Byte)

	for {
		err = p.Read()
		if err != nil {
			break
		}

		if p.Byte >= '0' && p.Byte <= '9' {
			sb.WriteByte(p.Byte)
		} else if p.Byte == '.' || p.Byte == 'e' || p.Byte == '-' || p.Byte == '+' {
			float = true
			sb.WriteByte(p.Byte)
		} else if p.Byte == 'E' {
			float = true
			sb.WriteByte('e')
		} else {
			break
		}
	}

	var (
		s  = sb.String()
		jv Value
	)

	if float {
		v, e := strconv.ParseFloat(s, 64)
		if e != nil {
			return nil, fmt.Errorf("parsing number at %d: %s", idx, e.Error())
		}
		jv = &floatValue{
			value: v,
			text:  s,
		}
	} else {
		v, e := strconv.Atoi(s)
		if e != nil {
			if strings.Contains(e.Error(), "value out of range") && v > 0 {
				var u uint64
				u, e = strconv.ParseUint(s, 10, 0)
				if e == nil {
					jv = &uintValue{
						value: u,
						text:  s,
					}
				}
			}
			if e != nil {
				return nil, fmt.Errorf("parsing number at %d: %s", idx, e.Error())
			}
		} else {
			jv = &intValue{
				value: v,
				text:  s,
			}
		}
	}

	if isWS(p.Byte) && err == nil {
		err = p.SkipWS()
	}

	return jv, err
}
