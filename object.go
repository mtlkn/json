package json

import (
	"strings"
)

type Property struct {
	Name  string
	Value *Value
}

type Object struct {
	Properties []*Property
	names      map[string]int
}

func New() *Object {
	return &Object{}
}

func (jo *Object) GetProperty(name string) (*Property, bool) {
	if jo == nil || len(jo.Properties) == 0 {
		return nil, false
	}

	jo.indexNames()
	i, ok := jo.names[name]
	if !ok {
		return nil, false
	}

	return jo.Properties[i], true
}

func (jo *Object) GetString(name string) (string, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return "", false
	}
	return jp.Value.GetString()
}

func (jo *Object) GetInt(name string) (int, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return 0, false
	}
	return jp.Value.GetInt()
}

func (jo *Object) GetUInt(name string) (uint, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return 0, false
	}
	return jp.Value.GetUInt()
}

func (jo *Object) GetFloat(name string) (float64, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return 0, false
	}
	return jp.Value.GetFloat()
}

func (jo *Object) GetBool(name string) (bool, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return false, false
	}
	return jp.Value.GetBool()
}

func (jo *Object) GetObject(name string) (*Object, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return nil, false
	}
	return jp.Value.GetObject()
}

func (jo *Object) GetArray(name string) (*Array, bool) {
	jp, ok := jo.GetProperty(name)
	if !ok {
		return nil, false
	}
	return jp.Value.GetArray()
}

func (jo *Object) GetStrings(name string) ([]string, bool) {
	ja, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return ja.GetStrings()
}

func (jo *Object) GetInts(name string) ([]int, bool) {
	ja, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return ja.GetInts()
}

func (jo *Object) GetFloats(name string) ([]float64, bool) {
	ja, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return ja.GetFloats()
}

func (jo *Object) GetObjects(name string) ([]*Object, bool) {
	ja, ok := jo.GetArray(name)
	if !ok {
		return nil, false
	}
	return ja.GetObjects()
}

func (jo *Object) Add(name string, value interface{}) *Object {
	jo.indexNames()

	v := newValue(value)

	i, ok := jo.names[name]
	if ok {
		jo.Properties[i].Value = v
	} else {
		jo.names[name] = len(jo.Properties)
		jo.Properties = append(jo.Properties, &Property{
			Name:  name,
			Value: v,
		})
	}

	return jo
}

func (jo *Object) Remove(name string) *Object {
	jo.indexNames()

	i, ok := jo.names[name]
	if !ok {
		return jo
	}

	names := make(map[string]int)
	var (
		jps   []*Property
		count int
	)

	for k, jp := range jo.Properties {
		if k == i {
			continue
		}

		names[jp.Name] = count
		jps = append(jps, jp)
		count++
	}

	jo.Properties = jps
	jo.names = names

	return jo
}

func (jo *Object) Merge(merge *Object) *Object {
	if merge == nil || len(merge.Properties) == 0 {
		return jo
	}

	for _, jp := range merge.Properties {
		v, err := jp.Value.GetValue()
		if err == nil {
			jo.Add(jp.Name, v)
		}
	}

	return jo
}

func (jo *Object) String() string {
	if jo == nil || len(jo.Properties) == 0 {
		return "{}"
	}

	var sb strings.Builder

	sb.WriteByte('{')

	for i, jp := range jo.Properties {
		if i > 0 {
			sb.WriteByte(',')
		}

		sb.WriteByte('"')
		sb.WriteString(jp.Name)
		sb.WriteByte('"')
		sb.WriteByte(':')
		sb.WriteString(jp.Value.String())
	}

	sb.WriteByte('}')

	return sb.String()
}

func (jo *Object) Bytes() []byte {
	return stringToBytes(jo.String())
}

func (jo *Object) Equals(other *Object) bool {
	if jo == nil && other != nil {
		return false
	}

	if jo != nil && other == nil {
		return false
	}

	if len(jo.Properties) != len(other.Properties) {
		return false
	}

	for _, jp := range jo.Properties {
		i, ok := other.names[jp.Name]
		if !ok {
			return false
		}

		if _, err := jp.Value.GetValue(); err != nil {
			return false
		}

		rv := other.Properties[i].Value

		if _, err := rv.GetValue(); err != nil {
			return false
		}

		if !jp.Value.Equals(rv) {
			return false
		}
	}

	return true
}

func (jo *Object) indexNames() {
	if jo.names != nil {
		return
	}

	jo.names = make(map[string]int)
	for i, jp := range jo.Properties {
		jo.names[jp.Name] = i
	}
}
