package json

import (
	"fmt"
)

type Property struct {
	Name   string
	Value  Value
	namep  []Parameter //name parameters
	valuep []Parameter //value parameters
}

func Field(name string, value Value) Property {
	return Property{
		Name:  name,
		Value: value,
	}
}

func (jp *Property) GetString() (v string, ok bool) {
	return StringValue(jp.Value)
}

func (jp *Property) GetStrings() ([]string, bool) {
	ja, ok := jp.GetArray()
	if !ok {
		return nil, false
	}

	return ja.GetStrings()
}

func (jp *Property) GetInt() (v int, ok bool) {
	return IntValue(jp.Value)
}

func (jp *Property) GetInts() ([]int, bool) {
	ja, ok := jp.GetArray()
	if !ok {
		return nil, false
	}

	return ja.GetInts()
}

func (jp *Property) GetFloat() (v float64, ok bool) {
	return FloatValue(jp.Value)
}

func (jp *Property) GetFloats() ([]float64, bool) {
	ja, ok := jp.GetArray()
	if !ok {
		return nil, false
	}

	return ja.GetFloats()
}

func (jp *Property) GetBool() (v bool, ok bool) {
	return BoolValue(jp.Value)
}

func (jp *Property) GetObject() (v *Object, ok bool) {
	return ObjectValue(jp.Value)
}

func (jp *Property) GetObjects() ([]*Object, bool) {
	ja, ok := jp.GetArray()
	if !ok {
		return nil, false
	}

	return ja.GetObjects()
}

func (jp *Property) GetArray() (v *Array, ok bool) {
	return ArrayValue(jp.Value)
}

//used when a first byte is '"'
func (p *byteParser) ParseProperty(parameterized bool) (*Property, error) {
	var (
		jp  Property
		err error
	)
	jp.Name, jp.namep, err = p.ParsePropertyName(parameterized)
	if err != nil {
		return &jp, err
	}

	if p.Byte != ':' {
		return nil, fmt.Errorf("parsing property at %d: expected [ : ], found %s", p.Index, string(p.Byte))
	}

	jp.Value, jp.valuep, err = p.ParseValue(parameterized)
	if err != nil {
		return &jp, fmt.Errorf("parsing \"%s\" property at %d: %s", jp.Name, p.Index, err.Error())
	}
	return &jp, nil
}

func (p *byteParser) ParsePropertyName(parameterized bool) (string, []Parameter, error) {
	var (
		idx    = p.Index
		params []Parameter
	)
	for {
		err := p.Read()
		if err != nil {
			return "", params, fmt.Errorf("parsing property name at %d: %s", idx, err.Error())
		}

		if p.Byte == '"' {
			name := string(p.Bytes[idx+1 : p.Index])
			return name, params, p.SkipWS()
		} else if p.Byte == '$' && parameterized {
			param, err := p.ParseParameter(idx)
			if param.Value != nil {
				params = append(params, param)
			}
			if err != nil {
				return "", params, fmt.Errorf("parsing property name at %d: %s", idx, err.Error())
			}
		}
	}
}
