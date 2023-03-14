package json

import (
	"strings"
)

type Property struct {
	Name  string
	Value *Value
}

type Object struct {
	Properites []*Property
	names      map[string]int
}

func New() *Object {
	return &Object{}
}

func (jo *Object) String() string {
	if jo == nil {
		return ""
	}

	if len(jo.Properites) == 0 {
		return "{}"
	}

	var sb strings.Builder

	sb.WriteByte('{')

	for i, jp := range jo.Properites {
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

func (jo *Object) GetProperty(name string) (*Property, bool) {
	if jo == nil || len(jo.Properites) == 0 {
		return nil, false
	}

	jo.indexNames()
	i, ok := jo.names[name]
	if !ok {
		return nil, false
	}

	return jo.Properites[i], true
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

func (jo *Object) Add(name string, value interface{}) *Object {
	jo.indexNames()

	v := newValue(value)

	i, ok := jo.names[name]
	if ok {
		jo.Properites[i].Value = v
	} else {
		jo.names[name] = len(jo.Properites)
		jo.Properites = append(jo.Properites, &Property{
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

	for k, jp := range jo.Properites {
		if k == i {
			continue
		}

		names[jp.Name] = count
		jps = append(jps, jp)
		count++
	}

	jo.Properites = jps
	jo.names = names

	return jo
}

func (jo *Object) indexNames() {
	if jo.names != nil {
		return
	}

	jo.names = make(map[string]int)
	for i, jp := range jo.Properites {
		jo.names[jp.Name] = i
	}
}
