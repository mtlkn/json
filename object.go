package json

import (
	"errors"
	"io"
	"strings"
)

type Object struct {
	Properties []*Property
	names      map[string]int
}

func NewObject(properties ...*Property) *Object {
	jo := &Object{
		names: make(map[string]int),
	}

	for _, p := range properties {
		i, ok := jo.names[p.Name()]
		if ok {
			jo.Properties[i] = p
			continue
		}

		jo.names[p.Name()] = len(jo.Properties)
		jo.Properties = append(jo.Properties, p)
	}

	return jo
}

// shortcut for NewObject
func O(properties ...*Property) *Object {
	return NewObject(properties...)
}

func (jo *Object) Get(name string) (*Value, bool) {
	i := jo.getPropertyIndex(name)
	if i == -1 {
		return nil, false
	}

	return jo.Properties[i].Value(), true
}

func (jo *Object) GetString(name string) (string, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return "", false
	}
	return v.String()
}

func (jo *Object) GetInt(name string) (int, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return 0, false
	}
	return v.Int()
}

func (jo *Object) GetUInt(name string) (uint, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return 0, false
	}
	return v.UInt()
}

func (jo *Object) GetFloat(name string) (float64, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return 0, false
	}
	return v.Float()
}

func (jo *Object) GetBool(name string) (bool, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return false, false
	}
	return v.Bool()
}

func (jo *Object) GetObject(name string) (*Object, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return nil, false
	}
	return v.Object()
}

func (jo *Object) GetArray(name string) (*Array, bool) {
	v, ok := jo.Get(name)
	if !ok {
		return nil, false
	}
	return v.Array()
}

func (jo *Object) GetStrings(name string) ([]string, bool) {
	a, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return a.GetStrings()
}

func (jo *Object) GetInts(name string) ([]int, bool) {
	a, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return a.GetInts()
}

func (jo *Object) GetFloats(name string) ([]float64, bool) {
	a, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return a.GetFloats()
}

func (jo *Object) GetObjects(name string) ([]*Object, bool) {
	a, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return a.GetObjects()
}

func (jo *Object) Set(name string, value interface{}) *Object {
	v := New(value)
	if v == nil {
		return jo
	}

	i := jo.getPropertyIndex(name)
	if i != -1 {
		jo.Properties[i].val = v
		return jo
	}

	jo.names[name] = len(jo.Properties)

	jo.Properties = append(jo.Properties, &Property{
		name: name,
		val:  v,
	})

	return jo
}

func (jo *Object) Remove(name string) *Object {
	r := jo.getPropertyIndex(name)
	if r == -1 {
		return jo
	}

	jo.names = make(map[string]int)

	var props []*Property

	for i, p := range jo.Properties {
		if i == r {
			continue
		}
		jo.names[p.Name()] = len(props)
		props = append(props, p)
	}

	jo.Properties = props

	return jo
}

func (jo *Object) Validate() error {
	for _, p := range jo.Properties {
		if p.Name() == "" {
			return errors.New("missing property name")
		}

		err := p.Value().Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (jo *Object) String() string {
	var sb strings.Builder

	sb.WriteByte('{')

	for i, p := range jo.Properties {
		if i > 0 {
			sb.WriteByte(',')
		}

		if len(p.n) > 0 {
			sb.Write(p.n)
		} else {
			sb.WriteByte('"')
			sb.WriteString(p.Name())
			sb.WriteByte('"')
		}

		sb.WriteByte(':')

		if len(p.v) > 0 {
			sb.Write(p.v)
			continue
		}

		v := p.Value()
		sb.WriteString(v.string())
	}

	sb.WriteByte('}')

	return sb.String()
}

func ParseObject(r io.Reader) (*Object, error) {
	rd, err := newReader(r)
	if err != nil {
		return nil, err
	}

	if !rd.SkipSpace() || rd.b != '{' {
		err = rd.EnsureJSON('{')
		if err != nil {
			return nil, err
		}
	}

	return rd.parseObject()
}

func (rd *reader) parseObject() (*Object, error) {
	jo := &Object{
		names: make(map[string]int),
	}

	for {
		p, err := rd.parseProperty()
		if err != nil {
			return nil, err
		}

		if p == nil {
			break
		}

		if _, ok := jo.names[p.Name()]; ok {
			return nil, errors.New("property name dublicates: " + p.Name())
		}

		jo.names[p.Name()] = len(jo.Properties)
		jo.Properties = append(jo.Properties, p)

		if rd.b == '}' {
			break
		}

		if rd.b != ',' {
			return nil, errors.New("missing property closing")
		}
	}

	return jo, nil
}

func (jo *Object) getPropertyIndex(name string) int {
	if jo == nil {
		return -1
	}

	// small array iteration is faster tham map lookup
	if len(jo.Properties) < 5 {
		for i, p := range jo.Properties {
			if p.Name() == name {
				return i
			}
		}
		return -1
	}

	i, ok := jo.names[name]
	if !ok {
		return -1
	}

	return i
}
