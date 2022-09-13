package json

import (
	"bytes"
	"strings"
	"testing"
)

func TestArray(t *testing.T) {
	t.Run("string array", func(t *testing.T) {
		ja := A("Alex", "Victoria", "Olga")
		vs, ok := ja.GetStrings()
		if !ok || len(vs) != 3 {
			t.Fail()
			return
		}
		for i, v := range []string{"Alex", "Victoria", "Olga"} {
			if vs[i] != v {
				t.Errorf("wrong value: expected \"%s\", got \"%s\"", v, vs[i])
				t.Fail()
			}
		}

		ja = A()
		vs, ok = ja.GetStrings()
		if !ok || len(vs) != 0 {
			t.Fail()
			return
		}

		ja = A(1)
		_, ok = ja.GetStrings()
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("int array", func(t *testing.T) {
		ja := A(1, 2, 3)
		vs, ok := ja.GetInts()
		if !ok || len(vs) != 3 {
			t.Fail()
			return
		}
		for i, v := range []int{1, 2, 3} {
			if vs[i] != v {
				t.Errorf("wrong value: expected \"%d\", got \"%d\"", v, vs[0])
				t.Fail()
			}
		}

		ja = A()
		vs, ok = ja.GetInts()
		if !ok || len(vs) != 0 {
			t.Fail()
			return
		}

		ja = A("YM")
		_, ok = ja.GetInts()
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("float array", func(t *testing.T) {
		ja := A(3.14, 0.2, -0.3)
		vs, ok := ja.GetFloats()
		if !ok || len(vs) != 3 {
			t.Fail()
			return
		}
		for i, v := range []float64{3.14, 0.2, -0.3} {
			if vs[i] != v {
				t.Errorf("wrong value: expected \"%f\", got \"%f\"", v, vs[0])
				t.Fail()
			}
		}

		ja = A()
		vs, ok = ja.GetFloats()
		if !ok || len(vs) != 0 {
			t.Fail()
			return
		}

		ja = A("YM")
		_, ok = ja.GetFloats()
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("object array", func(t *testing.T) {
		ja := A(O(P("name", "YM"), P("age", 27)), O(P("name", "SV"), P("age", 26)))
		vs, ok := ja.GetObjects()
		if !ok || len(vs) != 2 {
			t.Fail()
			return
		}
		for i, v := range []*Object{O(P("name", "YM"), P("age", 27)), O(P("name", "SV"), P("age", 26))} {
			{
				l, _ := vs[i].GetString("name")
				r, _ := v.GetString("name")
				if l != r {
					t.Errorf("wrong value: expected \"%s\", got \"%s\"", r, l)
					t.Fail()
				}
			}
			{
				l, _ := vs[i].GetInt("age")
				r, _ := v.GetInt("age")
				if l != r {
					t.Errorf("wrong value: expected \"%d\", got \"%d\"", r, l)
					t.Fail()
				}
			}
		}

		ja = A()
		vs, ok = ja.GetObjects()
		if !ok || len(vs) != 0 {
			t.Fail()
			return
		}

		ja = A(1)
		_, ok = ja.GetObjects()
		if ok {
			t.Fail()
			return
		}
	})

	t.Run("multi-type array", func(t *testing.T) {
		ja := A("YM", 2, 3.14, true, O(P("name", "SV"), P("age", 26)))
		if len(ja.Values) != 5 {
			t.Fail()
			return
		}

		s, ok := ja.Values[0].String()
		if !ok || s != "YM" {
			t.Fail()
			return
		}

		i, ok := ja.Values[1].Int()
		if !ok || i != 2 {
			t.Fail()
			return
		}

		f, ok := ja.Values[2].Float()
		if !ok || f != 3.14 {
			t.Fail()
			return
		}

		b, ok := ja.Values[3].Bool()
		if !ok || !b {
			t.Fail()
			return
		}

		o, ok := ja.Values[4].Object()
		if !ok {
			t.Fail()
			return
		}
		s, ok = o.GetString("name")
		if !ok || s != "SV" {
			t.Fail()
			return
		}
		i, ok = o.GetInt("age")
		if !ok || i != 26 {
			t.Fail()
		}
	})
}

func TestArrayParse(t *testing.T) {
	t.Run("parse string array", func(t *testing.T) {
		s := `[ "Alex", "Victoria", "Olga" ]`
		ja, err := ParseArray(strings.NewReader(s))
		if err != nil || len(ja.Values) != 3 {
			t.Error(err)
			t.Fail()
		}

		err = ja.Validate()
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		vs, ok := ja.GetStrings()
		if !ok || len(vs) != 3 {
			t.Fail()
			return
		}

		for i, v := range []string{"Alex", "Victoria", "Olga"} {
			if vs[i] != v {
				t.Errorf("wrong value: expected \"%s\", got \"%s\"", v, vs[i])
				t.Fail()
			}
		}
	})

	t.Run("parse object array", func(t *testing.T) {
		s := `[ {"name":"YM", "age" : 27 }, {"name":"SV","age":26} ]`
		ja, err := ParseArray(strings.NewReader(s))
		if err != nil || len(ja.Values) != 2 {
			t.Error(err)
			t.Fail()
		}

		vs, ok := ja.GetObjects()
		if !ok || len(vs) != 2 {
			t.Fail()
			return
		}
		for i, v := range []*Object{O(P("name", "YM"), P("age", 27)), O(P("name", "SV"), P("age", 26))} {
			{
				l, _ := vs[i].GetString("name")
				r, _ := v.GetString("name")
				if l != r {
					t.Errorf("wrong value: expected \"%s\", got \"%s\"", r, l)
					t.Fail()
				}
			}
			{
				l, _ := vs[i].GetInt("age")
				r, _ := v.GetInt("age")
				if l != r {
					t.Errorf("wrong value: expected \"%d\", got \"%d\"", r, l)
					t.Fail()
				}
			}
		}

		t.Run("array with nils", func(t *testing.T) {
			ja := A(nil)
			s := ja.String()
			if s != "[null]" {
				t.Fail()
			}
		})

		t.Run("array with bad values", func(t *testing.T) {
			ja := A(new(Value))
			err = ja.Validate()
			if err == nil {
				t.Fail()
			}
		})

	})

	t.Run("parse array array", func(t *testing.T) {
		s := `[ [1,2,3], ["YM", "SV" ]] `
		ja, err := ParseArray(strings.NewReader(s))
		if err != nil || len(ja.Values) != 2 {
			t.Error(err)
			t.Fail()
		}

		a, ok := ja.Values[0].Array()
		if !ok || len(a.Values) != 3 {
			t.Fail()
			return
		}
		is, ok := a.GetInts()
		if !ok || is[1] != 2 {
			t.Fail()
			return
		}

		a, ok = ja.Values[1].Array()
		if !ok || len(a.Values) != 2 {
			t.Fail()
			return
		}
		ss, ok := a.GetStrings()
		if !ok || ss[1] != "SV" {
			t.Fail()
			return
		}
	})

	t.Run("parse empty array", func(t *testing.T) {
		s := "[]"
		ja, err := ParseArray(strings.NewReader(s))
		if err != nil || len(ja.Values) != 0 {
			t.Error(err)
			t.Fail()
		}
	})

	t.Run("parse bad array", func(t *testing.T) {
		s := ""
		ja, err := ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		err = ja.Validate()
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		s = "3.14"
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "["
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "[ 1, 2"
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "[ 1, 2 "
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "[ 1 2 ]"
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "[ { \"name:2"
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}

		s = "[ [ 1,"
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}
	})

	t.Run("parse array with non-ASCII characters", func(t *testing.T) {
		bs := []byte("   [1,2]   ")
		bs[0] = 18
		bs[len(bs)-2] = 7
		ja, err := ParseArray(bytes.NewReader(bs))
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = ja.Validate()
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		is, ok := ja.GetInts()
		if !ok || len(is) != 2 || is[1] != 2 {
			t.Fail()
		}

		s := " x [1,2]"
		_, err = ParseArray(strings.NewReader(s))
		if err == nil {
			t.Fail()
		}
	})
}
