package json

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestObject(t *testing.T) {
	t.Run("string property", func(t *testing.T) {
		jo := O(P("name", "YM"))
		s, ok := jo.GetString("name")
		if !ok {
			t.Fail()
			return
		}
		if s != "YM" {
			t.Errorf("wrong value: expected \"%s\", got \"%s\"", "YM", s)
			t.Fail()
		}

		_, ok = jo.GetString("age")
		if ok {
			t.Fail()
			return
		}

		jo = O(P("names", A("YM", "SV")))
		ss, ok := jo.GetStrings("names")
		if !ok || len(ss) != 2 {
			t.Fail()
			return
		}
		if ss[1] != "SV" {
			t.Errorf("wrong value: expected \"%s\", got \"%s\"", "SV", ss[1])
			t.Fail()
		}

		_, ok = jo.GetStrings("age")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("int property", func(t *testing.T) {
		jo := O(P("age", 27))
		i, ok := jo.GetInt("age")
		if !ok {
			t.Fail()
			return
		}
		if i != 27 {
			t.Errorf("wrong value: expected \"%d\", got \"%d\"", 27, i)
			t.Fail()
		}

		_, ok = jo.GetInt("name")
		if ok {
			t.Fail()
			return
		}

		jo = O(P("ages", A(27, 43)))
		is, ok := jo.GetInts("ages")
		if !ok || len(is) != 2 {
			t.Fail()
			return
		}
		if is[1] != 43 {
			t.Errorf("wrong value: expected \"%d\", got \"%d\"", 43, is[1])
			t.Fail()
		}

		_, ok = jo.GetInts("name")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("float property", func(t *testing.T) {
		jo := O(P("pi", 3.14))
		f, ok := jo.GetFloat("pi")
		if !ok {
			t.Fail()
			return
		}
		if f != 3.14 {
			t.Errorf("wrong value: expected \"%f\", got \"%f\"", 3.14, f)
			t.Fail()
		}

		_, ok = jo.GetFloat("name")
		if ok {
			t.Fail()
			return
		}

		jo = O(P("geo", A(27.345, 43.876)))
		fs, ok := jo.GetFloats("geo")
		if !ok || len(fs) != 2 {
			t.Fail()
			return
		}
		if fs[1] != 43.876 {
			t.Errorf("wrong value: expected \"%f\", got \"%f\"", 43.876, fs[1])
			t.Fail()
		}

		_, ok = jo.GetFloats("name")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("bool property", func(t *testing.T) {
		jo := O(P("flag", true))
		b, ok := jo.GetBool("flag")
		if !ok {
			t.Fail()
			return
		}
		if !b {
			t.Errorf("wrong value: expected \"true\", got \"%v\"", b)
			t.Fail()
		}

		_, ok = jo.GetBool("name")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("uint property", func(t *testing.T) {
		var v uint = 9223372036854775808
		jo := O(P("big", v))
		ui, ok := jo.GetUInt("big")
		if !ok {
			t.Fail()
			return
		}
		if ui != v {
			t.Errorf("wrong value: expected \"%v\", got \"%v\"", v, ui)
			t.Fail()
		}

		_, ok = jo.GetUInt("name")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("object property", func(t *testing.T) {
		jo := O(P("org", O(P("name", "AP"))))
		o, ok := jo.GetObject("org")
		if !ok {
			t.Fail()
			return
		}
		s, ok := o.GetString("name")
		if !ok {
			t.Fail()
			return
		}
		if s != "AP" {
			t.Errorf("wrong value: expected \"%s\", got \"%s\"", "AP", s)
			t.Fail()
		}

		_, ok = jo.GetObject("name")
		if ok {
			t.Fail()
			return
		}

		jo = O(P("orgs", A(O(P("name", "AP")))))
		os, ok := jo.GetObjects("orgs")
		if !ok || len(os) != 1 {
			t.Fail()
			return
		}
		s, ok = os[0].GetString("name")
		if !ok {
			t.Fail()
			return
		}
		if s != "AP" {
			t.Errorf("wrong value: expected \"%s\", got \"%s\"", "AP", s)
			t.Fail()
		}

		_, ok = jo.GetObjects("name")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("array property", func(t *testing.T) {
		jo := O(P("names", A("YM", "SV")))
		a, ok := jo.GetArray("names")
		if !ok || len(a.Values) != 2 {
			t.Fail()
			return
		}
		ss, ok := a.GetStrings()
		if !ok || len(ss) != 2 {
			t.Fail()
			return
		}
		if ss[1] != "SV" {
			t.Errorf("wrong value: expected \"%s\", got \"%s\"", "SV", ss[1])
			t.Fail()
		}

		_, ok = jo.GetArray("name")
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("set property", func(t *testing.T) {
		jo := O().Set("name", "YM").Set("age", 27).Set("pi", 3.14).Set("cool", true)
		if len(jo.Properties) != 4 {
			t.Fail()
			return
		}

		s, ok := jo.GetString("name")
		if !ok || s != "YM" {
			t.Fail()
		}

		i, ok := jo.GetInt("age")
		if !ok || i != 27 {
			t.Fail()
		}

		jo.Set("age", 26)
		if len(jo.Properties) != 4 {
			t.Fail()
			return
		}
		i, ok = jo.GetInt("age")
		if !ok || i != 26 {
			t.Fail()
		}

		jo.Set("flag", true)
		jo.Set("active", false)
		if len(jo.Properties) != 6 {
			t.Fail()
			return
		}
		b, _ := jo.GetBool("flag")
		if !b {
			t.Fail()
		}
		_, ok = jo.GetBool("name")
		if ok {
			t.Fail()
		}

		jo.Set("null", nil)
		if len(jo.Properties) != 6 {
			t.Fail()
			return
		}
	})

	t.Run("nil object", func(t *testing.T) {
		var jo *Object
		_, ok := jo.Get("name")
		if ok {
			t.Fail()
		}
	})

	t.Run("new object", func(t *testing.T) {
		jo := O(P("name", "YM"), P("age", 27), P("pi", 3.14), P("cool", true), P("age", 26))
		if len(jo.Properties) != 4 {
			t.Fail()
			return
		}

		i, ok := jo.GetInt("age")
		if !ok || i != 26 {
			t.Fail()
		}
	})
}

func TestObjectParse(t *testing.T) {
	t.Run("parse object", func(t *testing.T) {
		s := `{"name": "Yuri Metelkin", "age": 27, "pi": 3.14, "cool": true, "org": {"name":"AP"}, "geo":[-0.1237, 3.214]}`
		jo, err := ParseObject(strings.NewReader(s))
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}

		err = jo.Validate()
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}

		fmt.Println(jo.String())

		s, ok := jo.GetString("name")
		if !ok || s != "Yuri Metelkin" {
			t.Fail()
		}

		i, ok := jo.GetInt("age")
		if !ok || i != 27 {
			t.Fail()
		}

		f, ok := jo.GetFloat("pi")
		if !ok || f != 3.14 {
			t.Fail()
		}

		b, ok := jo.GetBool("cool")
		if !ok || !b {
			t.Fail()
		}

		o, ok := jo.GetObject("org")
		if !ok {
			t.Fail()
			return
		}
		s, ok = o.GetString("name")
		if !ok || s != "AP" {
			t.Fail()
		}

		fs, ok := jo.GetFloats("geo")
		if !ok || len(fs) != 2 || fs[1] != 3.214 {
			t.Fail()
		}
	})

	t.Run("parse bad input object", func(t *testing.T) {
		s := `{"": "YM"}`
		jo, err := ParseObject(strings.NewReader(s))
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		err = jo.Validate()
		if err == nil {
			t.Fail()
		}

		s = ""
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = "3.14"
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = "{ "
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = "{ name: YM }"
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{"name" YM}`
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{"name": "YM" `
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{ "name": , "age": 27 }`
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{ "name": { name: "YM"} }`
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{ "name": [ "YM ] }`
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{ "name": "YM", "age": 27, "cool": true, "age": 26 }`
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
			return
		}

		s = `{ "name": YM }`
		jo, err = ParseObject(strings.NewReader(s))
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		err = jo.Validate()
		if err == nil {
			t.Fail()
		}

		s = "{ }"
		jo, err = ParseObject(strings.NewReader(s))
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		if len(jo.Properties) > 0 {
			t.Fail()
		}
	})

	t.Run("parse object with non-ASCII characters", func(t *testing.T) {
		bs := []byte("   { \"name\": \"YM\"}   ")
		bs[0] = 18
		bs[len(bs)-2] = 7
		jo, err := ParseObject(bytes.NewReader(bs))
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = jo.Validate()
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		s, ok := jo.GetString("name")
		if !ok || s != "YM" {
			t.Fail()
		}

		s = "x{ \"name\": \"YM\"}"
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "x"
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "  "
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "[1]"
		_, err = ParseObject(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		bs = []byte(" { \"name\": \"YM\"} x ")
		bs[0] = 18
		_, err = ParseObject(bytes.NewReader(bs))
		if err == nil {
			t.Fail()
		}
		_, err = ParseArray(bytes.NewReader(bs))
		if err == nil {
			t.Fail()
		}

		bs = []byte(" { \"name\": \"YM\"]")
		bs[0] = 18
		_, err = ParseObject(bytes.NewReader(bs))
		if err == nil {
			t.Fail()
		}
	})
}

func TestSetObject(t *testing.T) {
	jo := O(
		P("name", "Yuri Metelkin"),
		P("age", 26),
		P("pi", 3.14),
		P("cool", true),
		P("work", O(P("name", "AP"), P("title", "Developer"))),
		P("kids", A("Alex", "Victoria", "Olga")),
	)

	s, ok := jo.GetString("name")
	if !ok || s != "Yuri Metelkin" {
		t.Fail()
		return
	}

	i, ok := jo.GetInt("age")
	if !ok || i != 26 {
		t.Fail()
		return
	}

	f, ok := jo.GetFloat("pi")
	if !ok || f != 3.14 {
		t.Fail()
		return
	}

	b, ok := jo.GetBool("cool")
	if !ok || !b {
		t.Fail()
		return
	}

	o, ok := jo.GetObject("work")
	if !ok {
		t.Fail()
		return
	}
	s, ok = o.GetString("name")
	if !ok || s != "AP" {
		t.Fail()
		return
	}

	ss, ok := jo.GetStrings("kids")
	if !ok || ss[0] != "Alex" || ss[1] != "Victoria" || ss[2] != "Olga" {
		t.Fail()
		return
	}

	s = jo.String()
	fmt.Println(s)

	jo, err := ParseObject(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	s, ok = jo.GetString("name")
	if !ok || s != "Yuri Metelkin" {
		t.Fail()
		return
	}

	i, ok = jo.GetInt("age")
	if !ok || i != 26 {
		t.Fail()
		return
	}

	f, ok = jo.GetFloat("pi")
	if !ok || f != 3.14 {
		t.Fail()
		return
	}

	b, ok = jo.GetBool("cool")
	if !ok || !b {
		t.Fail()
		return
	}

	o, ok = jo.GetObject("work")
	if !ok {
		t.Fail()
		return
	}
	s, ok = o.GetString("name")
	if !ok || s != "AP" {
		t.Fail()
		return
	}

	ss, ok = jo.GetStrings("kids")
	if !ok || ss[0] != "Alex" || ss[1] != "Victoria" || ss[2] != "Olga" {
		t.Fail()
		return
	}

	fmt.Println(jo.String())

	jo.Remove("pi")
	_, ok = jo.GetBool("pi")
	if ok {
		t.Fail()
	}
	fmt.Println(jo.String())

	jo.Remove("pi")
	fmt.Println(jo.String())
}
