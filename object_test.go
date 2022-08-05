package json

import (
	"fmt"
	"os"
	"testing"
)

func TestObject(t *testing.T) {
	jo := New(Field("id", Int(1)))
	jo.addProperty(&Property{
		Name:  "id",
		Value: String("1"),
	})
	if s, _ := jo.GetString("id"); s != "1" {
		t.Fail()
		t.Error("failed to update property")
	}

	jo.Add("name", String("YM"))
	jo.Remove("id")
	if jo.Properties[0].Name != "name" {
		t.Fail()
		t.Error("failed to remove property")
	}

	jo.Remove("age")
	if jo.Properties[0].Name != "name" {
		t.Fail()
		t.Error("failed to remove property")
	}

	jo.Remove("name")
	if len(jo.Properties) != 0 {
		t.Fail()
		t.Error("failed to remove property")
	}

	jo.Remove("name")
	if len(jo.Properties) != 0 {
		t.Fail()
		t.Error("failed to remove property")
	}

	if _, ok := jo.GetProperty("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetString("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetStrings("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetInt("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetInts("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetFloat("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetFloats("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetBool("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetObject("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetObjects("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	if _, ok := jo.GetArray("id"); ok {
		t.Fail()
		t.Error("failed to fail to get nonexisting property")
	}

	jo.Add("a", NewStringArray([]string{"a"}))
	if _, ok := jo.GetStrings("a"); !ok {
		t.Fail()
		t.Error("failed to get array property")
	}

	jo.Add("a", NewIntArray([]int{1}))
	if _, ok := jo.GetInts("a"); !ok {
		t.Fail()
		t.Error("failed to get array property")
	}

	jo.Add("a", NewFloatArray([]float64{3.14}))
	if _, ok := jo.GetFloats("a"); !ok {
		t.Fail()
		t.Error("failed to get array property")
	}

	jo.Add("a", NewObjectArray([]*Object{New()}))
	if _, ok := jo.GetObjects("a"); !ok {
		t.Fail()
		t.Error("failed to get array property")
	}
	if _, ok := jo.GetArray("a"); !ok {
		t.Fail()
		t.Error("failed to get array property")
	}

	jo.Add("b", Bool(true))
	if _, ok := jo.GetBool("b"); !ok {
		t.Fail()
		t.Error("failed to get bool property")
	}

	jo.Add("o", New())
	if _, ok := jo.GetObject("o"); !ok {
		t.Fail()
		t.Error("failed to get object property")
	}

	jo = nil
	jo = jo.Copy()
	if jo != nil {
		t.Fail()
		t.Error("failed to copy nil object")
	}

	jo = New()
	jo = jo.Copy()
	if jo == nil {
		t.Fail()
		t.Error("failed to copy empty object")
	}

	jo, _ = ParseObjectWithParameters([]byte(`{"${name}":"${value}"}`))
	jo = jo.Copy()
	jp := jo.Properties[0]
	if len(jp.namep) == 0 || jp.namep[0].Name != "name" || len(jp.valuep) == 0 || jp.valuep[0].Name != "value" {
		t.Fail()
		t.Error("failed to copy parameterized object")
	}

	var l, r *Object
	if ok, _ := l.Equals(r); !ok {
		t.Fail()
		t.Error("failed to compare nil objects")
	}

	l = New()
	r = New(Field("id", Int(1)))
	if ok, _ := l.Equals(r); ok {
		t.Fail()
		t.Error("failed to fail to compare different objects")
	}

	jo, _ = ParseObjectWithParameters([]byte(`{"${name}":"${value}"}`))
	params := New(Field("id", String("x")))
	o := jo.SetParameters(params)
	if len(o.Properties) > 0 {
		t.Fail()
		t.Error("failed to fail set missing name parameter")
	}

	params.Add("name", String("title"))
	o = jo.SetParameters(params)
	if len(o.Properties) > 0 {
		t.Fail()
		t.Error("failed to fail set missing value parameters")
	}

	if jo.Value() != jo {
		t.Fail()
		t.Error("failed to compare object value to itself")
	}

	jo = nil
	if jo.String() != "{}" {
		t.Fail()
		t.Error("failed to string nil object")
	}

	jo = New()
	if jo.String() != "{}" {
		t.Fail()
		t.Error("failed to string empty object")
	}
}

func TestObjectParse(t *testing.T) {
	s := `{ "text": "abc", "number": 3.14, "flag": true, "array": [ 1, 2, 3 ], "object": { "a": "b" }}`
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '{' {
		t.Error("Failed to parse {")
	}
	v, err := p.ParseObject(false)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(v.String())

	s = `{"x": "Arts &amp; Entertainment; a &lt; b or c &gt; d; YM & &Co"}`
	p = newParser([]byte(s))
	p.SkipWS()
	v, err = p.ParseObject(false)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v.String())

	if _, err := parseObject([]byte("[]"), false, false); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad json")
	}

	if _, err := parseObject([]byte("{"), false, false); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad json")
	}

	if _, err := parseObject([]byte("{f"), false, false); err == nil {
		t.Fail()
		t.Error("failed to fail to parse bad json")
	}

	if _, err := parseObject([]byte("{}"), false, false); err != nil {
		t.Fail()
		t.Error("failed to parse empty json")
	}

	if _, err := parseObject([]byte("{\"f\" : { \"age\": 27 }  }"), false, false); err != nil {
		t.Fail()
		t.Error("failed to parse json")
	}
}

func TestEnsureJSON(t *testing.T) {
	s := `  ï » ¿{ "text": "abc"} ï » ¿ `
	jo, err := ParseObjectSafe([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(jo.String())

	s = ` ï » ¿<text>abc</text>`
	_, err = ParseObjectSafe([]byte(s))
	if err == nil {
		t.Error("Must not parse it")
	}
	fmt.Println(err.Error())

	s = `{ "text": "abc"]`
	_, err = ParseObjectSafe([]byte(s))
	if err == nil {
		t.Error("Must not parse it")
	}
	fmt.Println(err.Error())

	s = `  ï » ¿[{ "text": "abc"}] ï » ¿ `
	ja, err := ParseArraySafe([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(ja.String())
}

func TestObjectCopy(t *testing.T) {
	jo := New(
		Field("name", String("YM")),
		Field("null", Null()),
		Field("empty", String("")),
	)
	fmt.Println(jo.String())
	copy := jo.Copy()
	fmt.Println(copy.String())
}

func TestObjectPointers(t *testing.T) {
	jo := New(Field("name", String("YM")))
	fmt.Println(jo.String())
	jo.Add("person", jo)
	fmt.Println(jo.String())
}

func TestGraph(t *testing.T) {
	data, _ := os.ReadFile("testdata/graph.json")
	jo, _ := ParseObject(data)
	ja, _ := jo.GetObjects("vertices")
	vertices := make(map[int]graphPerson)
	for i, v := range ja {
		name, _ := v.GetString("term")
		vertices[i] = graphPerson{
			Name: name,
		}
	}

	ja, _ = jo.GetObjects("connections")
	for _, o := range ja {
		source, _ := o.GetInt("source")
		target, _ := o.GetInt("target")
		weight, _ := o.GetFloat("weight")
		count, _ := o.GetInt("doc_count")
		v := vertices[source]
		c := vertices[target]
		v.Connections = append(v.Connections, graphConnection{
			Name:   c.Name,
			Weight: weight,
			Count:  count,
		})
		vertices[source] = v
	}

	for _, p := range vertices {
		if len(p.Connections) == 0 {
			continue
		}

		fmt.Println(p.Name)
		for _, c := range p.Connections {
			fmt.Printf("\t%s\n", c.Name)
		}
		fmt.Println()
	}
}

type graphPerson struct {
	Name        string
	Connections []graphConnection
}

type graphConnection struct {
	Name   string
	Weight float64
	Count  int
}
