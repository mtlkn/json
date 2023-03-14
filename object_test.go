package json

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestObject(t *testing.T) {
	t.Run("parse simple object", func(t *testing.T) {
		bs, _ := os.ReadFile("testdata/simple.json")
		jo, err := ParseObject(bs)
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}

		if len(jo.Properites) != 11 {
			t.Fail()
		}

		for _, jp := range jo.Properites {
			v, err := jp.Value.GetValue()
			if err != nil {
				t.Fail()
				t.Error(err)
			}
			fmt.Println(jp.Name, jp.Value.Type.String(), v)
		}

		fmt.Println(jo.String())

		v, err := Parse(bs)
		if err != nil {
			t.Fail()
			t.Error()
			return
		}

		if v.Type != OBJECT {
			t.Fail()
		}
	})

	t.Run("parse complex object", func(t *testing.T) {
		bs, _ := os.ReadFile("testdata/complex.json")
		jo, err := ParseObject(bs)
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}

		for _, jp := range jo.Properites {
			v, err := jp.Value.GetValue()
			if err != nil {
				t.Fail()
				t.Error(err)
			}
			fmt.Println(jp.Name, jp.Value.Type.String(), v)
		}

		fmt.Println(jo.String())
	})

	t.Run("object getters", func(t *testing.T) {
		bs, _ := os.ReadFile("testdata/simple.json")
		jo, _ := ParseObject(bs)

		s, ok := jo.GetString("name")
		if !ok || s != "YM" {
			t.Fail()
		}

		_, ok = jo.GetString("age")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetString("xyz")
		if ok {
			t.Fail()
		}

		d, ok := jo.GetInt("age")
		if !ok || d != 27 {
			t.Fail()
		}

		_, ok = jo.GetInt("name")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetInt("xyz")
		if ok {
			t.Fail()
		}

		f, ok := jo.GetFloat("comp")
		if !ok || f != 250000.99 {
			t.Fail()
		}

		_, ok = jo.GetFloat("age")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetFloat("xyz")
		if ok {
			t.Fail()
		}

		b, ok := jo.GetBool("true")
		if !ok || !b {
			t.Fail()
		}

		b, ok = jo.GetBool("false")
		if !ok || b {
			t.Fail()
		}

		_, ok = jo.GetBool("age")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetBool("xyz")
		if ok {
			t.Fail()
		}

		o, ok := jo.GetObject("obj")
		if !ok || o == nil || len(o.Properites) != 5 {
			t.Fail()
		}

		_, ok = jo.GetObject("age")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetObject("xyz")
		if ok {
			t.Fail()
		}

		a, ok := jo.GetArray("array")
		if !ok || o == nil || len(a.Values) != 5 {
			t.Fail()
		}

		_, ok = jo.GetArray("age")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetArray("xyz")
		if ok {
			t.Fail()
		}

		bs, _ = os.ReadFile("testdata/complex.json")
		jo, _ = ParseObject(bs)

		u, ok := jo.GetUInt("uint")
		if !ok || u != 9223372036854775888 {
			t.Fail()
		}

		_, ok = jo.GetUInt("age")
		if ok {
			t.Fail()
		}

		_, ok = jo.GetUInt("xyz")
		if ok {
			t.Fail()
		}
	})

	t.Run("add properties", func(t *testing.T) {
		jo := New()
		jo.Add("name", "YM").Add("age", 27).Add("ap", false)
		jo.Add("test", []float32{3.14})
		jo.Add("ap", true)
		if len(jo.Properites) != 4 {
			t.Fail()
		}
		fmt.Println(jo.String())

		jo.Remove("xyz")
		if len(jo.Properites) != 4 {
			t.Fail()
		}

		jo.Remove("ap")
		if len(jo.Properites) != 3 {
			t.Fail()
		}
		fmt.Println(jo.String())

		jo.Remove("name")
		if len(jo.Properites) != 2 {
			t.Fail()
		}
		fmt.Println(jo.String())

		jo.Remove("test")
		if len(jo.Properites) != 1 {
			t.Fail()
		}
		fmt.Println(jo.String())
	})

	t.Run("errors", func(t *testing.T) {
		var jo *Object
		s := jo.String()
		if s != "" {
			t.Fail()
		}

		jo = new(Object)
		s = jo.String()
		if s != "{}" {
			t.Fail()
		}

		_, ok := jo.GetString("name")
		if ok {
			t.Fail()
		}
	})
}

func TestProperty(t *testing.T) {
	t.Run("parse string property", func(t *testing.T) {
		s := `{ "name": "YM" }`
		bs := []byte(s)
		last := len(bs)
		jp, end, err := parseProperty(bs, 1, last)
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}

		if end != last-1 {
			t.Fail()
		}

		if jp.Name != "name" {
			t.Fail()
		}

		if string(jp.Value.data) != `"YM"` {
			t.Fail()
		}
	})
}

func BenchmarkObjectParsers(b *testing.B) {
	bs, _ := os.ReadFile("testdata/simple.json")
	obj := make(map[string]interface{})

	b.Run("go parser", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			json.Unmarshal(bs, &obj)
		}
	})

	// b.Run("js parser", func(b *testing.B) {
	// 	for n := 0; n < b.N; n++ {
	// 		js.ParseObject(bs)
	// 	}
	// })

	b.Run("this parser", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			ParseObject(bs)
		}
	})
}

func BenchmarkLargeObjectParsers(b *testing.B) {
	bs, _ := os.ReadFile("testdata/appl.json")
	obj := make(map[string]interface{})

	b.Run("go parser", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			json.Unmarshal(bs, &obj)
		}
	})

	// b.Run("js parser", func(b *testing.B) {
	// 	for n := 0; n < b.N; n++ {
	// 		js.ParseObject(bs)
	// 	}
	// })

	b.Run("this parser", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			ParseObject(bs)
		}
	})
}
