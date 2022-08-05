package json

import (
	"fmt"
	"strings"
)

type Object struct {
	Properties []*Property
	fields     map[string]int
	params     bool
	text       string
}

func New(props ...Property) *Object {
	var jo Object
	if len(props) > 0 {
		for _, prop := range props {
			jo.Add(prop.Name, prop.Value)
		}
	}
	return &jo
}

func ParseObject(data []byte) (*Object, error) {
	return parseObject(data, false, false)
}

//ParseObjectSafe parses JSON object ignoring prefix non-ASCII characters
func ParseObjectSafe(data []byte) (*Object, error) {
	return parseObject(data, false, true)
}

//ParseObjectWithParameters parses parameterized JSON object
func ParseObjectWithParameters(data []byte) (*Object, error) {
	return parseObject(data, true, false)
}

func parseObject(data []byte, parameterized bool, safe bool) (*Object, error) {
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
	if p.Byte != '{' {
		return nil, fmt.Errorf("parsing JSON object; expect '{', found '%s'", string(p.Byte))
	}
	return p.ParseObject(parameterized)
}

func (jo *Object) Add(field string, value Value) {
	if value.Type() == ObjectType && jo == value {
		value, _ = copyValue(value)
	}

	jo.addProperty(&Property{
		Name:  field,
		Value: value,
	})
}

func (jo *Object) addProperty(jp *Property) {
	jo.text = ""

	if jo.fields == nil {
		jo.fields = make(map[string]int)
	}

	idx, ok := jo.fields[jp.Name]
	if ok {
		jo.Properties[idx] = jp
		return
	}

	jo.fields[jp.Name] = len(jo.Properties)
	jo.Properties = append(jo.Properties, jp)
}

func (jo *Object) Remove(field string) {
	sz := len(jo.fields)
	if sz == 0 {
		return
	}

	idx, ok := jo.fields[field]
	if !ok {
		return
	}

	delete(jo.fields, field)

	jo.text = ""

	props := make([]*Property, len(jo.fields))

	for i := 0; i < len(jo.Properties); i++ {
		if i == idx {
			continue
		}

		jp := jo.Properties[i]

		k := i
		if i > idx {
			k--
		}

		props[k] = jp
		jo.fields[jp.Name] = k
	}

	jo.Properties = props
}

func (jo *Object) GetProperty(field string) (*Property, bool) {
	if len(jo.fields) == 0 {
		return nil, false
	}

	idx, ok := jo.fields[field]
	if !ok {
		return nil, false
	}

	return jo.Properties[idx], true
}

func (jo *Object) GetValue(field string) (Value, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.Value, true
}

func (jo *Object) GetString(field string) (string, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return "", false
	}
	return jp.GetString()
}

func (jo *Object) GetStrings(field string) ([]string, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetStrings()
}

func (jo *Object) GetInt(field string) (int, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return 0, false
	}
	return jp.GetInt()
}

func (jo *Object) GetInts(field string) ([]int, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetInts()
}

func (jo *Object) GetFloat(field string) (float64, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return 0, false
	}
	return jp.GetFloat()
}

func (jo *Object) GetFloats(field string) ([]float64, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetFloats()
}

func (jo *Object) GetBool(field string) (bool, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return false, false
	}
	return jp.GetBool()
}

func (jo *Object) GetObject(field string) (*Object, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetObject()
}

func (jo *Object) GetObjects(field string) ([]*Object, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetObjects()
}

func (jo *Object) GetArray(field string) (*Array, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetArray()
}

func (jo *Object) Copy() *Object {
	if jo == nil {
		return nil
	}

	if len(jo.Properties) == 0 {
		return new(Object)
	}

	copy := Object{
		params: jo.params,
		fields: make(map[string]int),
	}

	for _, jp := range jo.Properties {
		v, ok := copyValue(jp.Value)
		if ok {
			copy.fields[jp.Name] = len(copy.Properties)

			cp := &Property{
				Name:  jp.Name,
				Value: v,
			}

			if len(jp.namep) > 0 {
				cp.namep = append(cp.namep, jp.namep...)
			}

			if len(jp.valuep) > 0 {
				cp.valuep = append(cp.valuep, jp.valuep...)
			}

			copy.Properties = append(copy.Properties, cp)
		}
	}

	return &copy
}

func (jo *Object) IsEmpty() bool {
	return jo == nil || len(jo.Properties) == 0
}

func (jo *Object) Equals(other *Object) (bool, error) {
	var (
		left  = jo
		right = other
	)

	if left == nil {
		left = New()
	}

	if right == nil {
		right = New()
	}

	for f := range left.fields {
		if _, ok := right.fields[f]; !ok {
			v, _ := left.GetValue(f)
			if !v.IsEmpty() {
				return false, fmt.Errorf("extra property: %s", f)
			}
		}
	}

	for f := range right.fields {
		if _, ok := left.fields[f]; !ok {
			v, _ := right.GetValue(f)
			if !v.IsEmpty() {
				return false, fmt.Errorf("missing property: %s", f)
			}
		}
	}

	for _, l := range left.Properties {
		for _, r := range right.Properties {
			if r.Name == l.Name {
				err := compareValues(l.Value, r.Value)
				if err != nil {
					return false, fmt.Errorf("mismatch property %s: %s", l.Name, err.Error())
				}
				break
			}
		}
	}

	return true, nil
}

//SetParameters replaces parameter placeholders with values
func (jo *Object) SetParameters(params *Object) *Object {
	var set Object

	for _, jp := range jo.Properties {
		var (
			name  = jp.Name
			value = jp.Value
		)
		if len(jp.namep) > 0 {
			name, _ = setStringParameters(fmt.Sprintf("\"%s\"", jp.Name), jp.namep, params)
			if len(name) < 3 {
				continue
			}
			name = string(name[1 : len(name)-1])
		}

		value = setValueParameters(value, jp.valuep, params)
		if value == nil || value.IsEmpty() {
			continue
		}
		set.Add(name, value)
	}

	return &set
}

//GetParameters retrieves paramaters from Object
func (jo *Object) GetParameters() []Parameter {
	var (
		params []Parameter
	)

	for _, jp := range jo.Properties {
		if len(jp.namep) > 0 {
			params = append(params, jp.namep...)
		}

		switch jp.Value.Type() {
		case StringType:
			if len(jp.valuep) > 0 {
				params = append(params, jp.valuep...)
			}
		case ObjectType:
			o, ok := jp.Value.(*Object)
			if ok {
				params = append(params, o.GetParameters()...)
			}
		case ArrayType:
			a, ok := jp.Value.(*Array)
			if ok {
				params = append(params, a.GetParameters()...)
			}
		}
	}

	return params
}

func (jo *Object) Value() interface{} {
	return jo
}

func (jo *Object) Type() ValueType {
	return ObjectType
}

func (jo *Object) String() string {
	if jo == nil {
		return "{}"
	}

	if jo.text == "" {
		sz := len(jo.Properties)
		if sz == 0 {
			jo.text = "{}"
		} else {
			values := make([]string, sz)
			for i, jp := range jo.Properties {
				values[i] = fmt.Sprintf("\"%s\":%s", jp.Name, jp.Value.String())
			}
			jo.text = fmt.Sprintf("{%s}", strings.Join(values, ","))
		}
	}
	return jo.text
}

//used when a first byte is '{'
func (p *byteParser) ParseObject(parameterized bool) (*Object, error) {
	var (
		i      int
		params bool
		props  []*Property
		fields = make(map[string]int)
	)

	for {
		err := p.SkipWS()
		if err != nil {
			return nil, err
		}
		if p.Byte != '"' {
			if p.Byte == '}' {
				break
			}
			return nil, fmt.Errorf("parsing object at %d: expected [ \" ], found %s", p.Index, string(p.Byte))
		}
		jp, err := p.ParseProperty(parameterized)
		if err != nil {
			return nil, err
		}

		if len(jp.valuep) > 0 || len(jp.namep) > 0 {
			params = true
		}

		end := p.Byte == '}'
		if end || p.Byte == ',' {
			fields[jp.Name] = i
			props = append(props, jp)
			i++

			if end {
				break
			}
			continue
		}

		return nil, fmt.Errorf("parsing object at %d: expected [ , } ], found %s", p.Index, string(p.Byte))
	}

	jo := &Object{
		Properties: props,
		fields:     fields,
		params:     params,
	}

	err := p.SkipWS()
	if err == errEOF {
		err = nil
	}

	return jo, err
}
