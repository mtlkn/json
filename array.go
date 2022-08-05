package json

import (
	"fmt"
	"strings"
)

type Array struct {
	Values []Value
	text   string
	params map[int][]Parameter
}

func ParseArray(data []byte) (*Array, error) {
	return parseArray(data, false, false)
}

func ParseArrayWithParameters(data []byte) (*Array, error) {
	return parseArray(data, true, false)
}

func ParseArraySafe(data []byte) (*Array, error) {
	return parseArray(data, false, true)
}

func parseArray(data []byte, parameterized bool, safe bool) (*Array, error) {
	var (
		p   = newParser(data)
		err error
	)

	if safe {
		err = p.EnsureJSON()
	} else {
		err = p.SkipWS()
	}
	if err != nil {
		return nil, err
	}
	if p.Byte != '[' {
		return nil, fmt.Errorf("parsing JSON array; expect '[', found '%s'", string(p.Byte))
	}
	return p.ParseArray(parameterized)
}

func NewArray(vs ...Value) *Array {
	return &Array{
		Values: vs,
	}
}

func NewStringArray(vs []string) *Array {
	var ja Array
	if len(vs) > 0 {
		ja.Values = make([]Value, len(vs))
		for i, v := range vs {
			ja.Values[i] = &stringValue{
				value: v,
			}
		}
	}
	return &ja
}

func NewIntArray(vs []int) *Array {
	var ja Array
	if len(vs) > 0 {
		ja.Values = make([]Value, len(vs))
		for i, v := range vs {
			ja.Values[i] = &intValue{
				value: v,
			}
		}
	}
	return &ja
}

func NewFloatArray(vs []float64) *Array {
	var ja Array
	if len(vs) > 0 {
		ja.Values = make([]Value, len(vs))
		for i, v := range vs {
			ja.Values[i] = &floatValue{
				value: v,
			}
		}
	}
	return &ja
}

func NewObjectArray(vs []*Object) *Array {
	var ja Array
	if len(vs) > 0 {
		ja.Values = make([]Value, len(vs))
		for i, v := range vs {
			ja.Values[i] = v
		}
	}
	return &ja
}

func (ja *Array) AddString(v string) {
	ja.Values = append(ja.Values, &stringValue{
		value: v,
	})
}

func (ja *Array) AddInt(v int) {
	ja.Values = append(ja.Values, &intValue{
		value: v,
	})
}

func (ja *Array) AddFloat(v float64) {
	ja.Values = append(ja.Values, &floatValue{
		value: v,
	})
}

func (ja *Array) AddObject(v *Object) {
	ja.Values = append(ja.Values, v)
}

func (ja *Array) GetStrings() ([]string, bool) {
	vs := make([]string, len(ja.Values))
	for i, jv := range ja.Values {
		v, ok := StringValue(jv)
		if !ok {
			return vs, false
		}
		vs[i] = v
	}
	return vs, true
}

func (ja *Array) GetInts() ([]int, bool) {
	vs := make([]int, len(ja.Values))
	for i, jv := range ja.Values {
		v, ok := IntValue(jv)
		if !ok {
			return vs, false
		}
		vs[i] = v
	}
	return vs, true
}

func (ja *Array) GetFloats() ([]float64, bool) {
	vs := make([]float64, len(ja.Values))
	for i, jv := range ja.Values {
		v, ok := FloatValue(jv)
		if !ok {
			return vs, false
		}
		vs[i] = v
	}
	return vs, true
}

func (ja *Array) GetObjects() ([]*Object, bool) {
	vs := make([]*Object, len(ja.Values))
	for i, jv := range ja.Values {
		v, ok := ObjectValue(jv)
		if !ok {
			return vs, false
		}
		vs[i] = v
	}
	return vs, true
}

func (ja *Array) Copy() *Array {
	if ja == nil {
		return nil
	}
	if len(ja.Values) == 0 {
		return &Array{}
	}

	var copy Array
	for _, jv := range ja.Values {
		v, ok := copyValue(jv)
		if ok {
			copy.Values = append(copy.Values, v)
		}
	}

	if len(ja.params) > 0 {
		copy.params = make(map[int][]Parameter)
		for i, v := range ja.params {
			copy.params[i] = v
		}
	}

	return &copy
}

func (ja *Array) Equals(other *Array) (bool, error) {
	var (
		left  = ja
		right = other
	)

	if left == nil {
		left = new(Array)
	}

	if right == nil {
		right = new(Array)
	}

	if len(left.Values) != len(right.Values) {
		return false, fmt.Errorf("mismatch array size: %d != %d", len(left.Values), len(right.Values))
	}

	for i, l := range left.Values {
		var err error
		for j, r := range right.Values {
			e := compareValues(l, r)
			if e == nil {
				err = nil
				break
			}
			if j == i {
				err = e
			}
		}
		if err != nil {
			return false, fmt.Errorf("value missing at [%d]: %s", i, err.Error())
		}
	}

	return true, nil
}

func (ja *Array) SetParameters(params *Object) *Array {
	var set Array

	for i, v := range ja.Values {
		pms := ja.params[i]
		v = setValueParameters(v, pms, params)
		if v == nil || v.IsEmpty() {
			continue
		}

		set.Values = append(set.Values, v)
	}

	return &set
}

func (ja *Array) GetParameters() []Parameter {
	var (
		params []Parameter
	)

	for i, v := range ja.Values {
		switch v.Type() {
		case StringType:
			pms, ok := ja.params[i]
			if ok && len(pms) > 0 {
				params = append(params, pms...)
			}
		case ObjectType:
			o, ok := v.(*Object)
			if ok {
				params = append(params, o.GetParameters()...)
			}
		case ArrayType:
			a, ok := v.(*Array)
			if ok {
				params = append(params, a.GetParameters()...)
			}
		}
	}

	return params
}

func (ja *Array) Value() interface{} {
	return ja.Values
}

func (ja *Array) Type() ValueType {
	return ArrayType
}

func (ja *Array) String() string {
	if ja == nil {
		return "[]"
	}

	if ja.text == "" {
		sz := len(ja.Values)
		if sz == 0 {
			ja.text = "[]"
		} else {
			values := make([]string, sz)
			for i, v := range ja.Values {
				values[i] = v.String()
			}
			ja.text = fmt.Sprintf("[%s]", strings.Join(values, ","))
		}
	}
	return ja.text
}

func (ja *Array) IsEmpty() bool {
	return ja == nil || len(ja.Values) == 0
}

//used when a first byte is '['
func (p *byteParser) ParseArray(parameterized bool) (*Array, error) {
	var (
		params map[int][]Parameter
		values []Value
	)
	for {
		idx := p.Index
		v, pms, err := p.ParseValue(parameterized)
		if err != nil {
			if p.Byte == ']' {
				break
			}
			return nil, fmt.Errorf("parsing array at %d: %s", idx, err.Error())
		}

		if len(pms) > 0 {
			if params == nil {
				params = make(map[int][]Parameter)
			}
			params[len(values)] = pms
		}

		if p.Byte == ',' {
			values = append(values, v)
			continue
		}

		if p.Byte == ']' {
			values = append(values, v)
			break
		}

	}

	err := p.SkipWS()
	if err == errEOF {
		err = nil
	}

	return &Array{
		Values: values,
		params: params,
	}, err
}
