package json

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {
	t.Run("simple array", func(t *testing.T) {
		s := `["abc",123,true, false, { "name": "YM" } ]`
		bs := []byte(s)

		ja, err := ParseArray(bs)
		if err != nil {
			t.Fail()
			t.Error()
			return
		}

		if len(ja.Values) != 5 {
			t.Fail()
		}

		for i, jv := range ja.Values {
			fmt.Println(i+1, jv.Type, string(jv.data))
		}

		v, err := Parse(bs)
		if err != nil {
			t.Fail()
			t.Error()
			return
		}

		if v.Type != ARRAY {
			t.Fail()
		}

		if _, err := ParseArrayString("[]"); err != nil {
			t.Fail()
			t.Error()
		}

		if _, err := ParseArrayString(" [ ] "); err != nil {
			t.Fail()
			t.Error()
		}

		if _, err := ParseObjectString(`{ "x": []}`); err != nil {
			t.Fail()
			t.Error()
		}
	})

	t.Run("create array", func(t *testing.T) {
		ja := NewArray([]string{"abc", "xyz"})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != `"abc"` {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]int{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]float64{3.14, 0.2e-2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "3.14" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]float32{3.14, 0.2e-2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "3.14" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]bool{true, false})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "true" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]uint{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]int64{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]int8{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]int16{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]int32{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]uint64{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]uint8{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]uint16{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]uint32{1, 2})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray("xyz")
		if ja == nil || len(ja.Values) != 1 || ja.Values[0].String() != `"xyz"` {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray(1)
		if ja == nil || len(ja.Values) != 1 || ja.Values[0].String() != "1" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray(3.14)
		if ja == nil || len(ja.Values) != 1 || ja.Values[0].String() != "3.14" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray(true)
		if ja == nil || len(ja.Values) != 1 || ja.Values[0].String() != "true" {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]*Object{New().Add("name", "YM"), New().Add("age", 27)})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].Type != OBJECT {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray([]*Array{NewArray("xyz"), NewArray([]int{1, 2})})
		if ja == nil || len(ja.Values) != 2 || ja.Values[0].Type != ARRAY {
			t.Fail()
		}
		fmt.Println(ja.String())

		ja = NewArray(nil)
		if ja != nil {
			t.Fail()
		}

		ja = NewArray([]*Property{{Name: "name"}})
		if ja != nil {
			t.Fail()
		}
	})

	t.Run("getters", func(t *testing.T) {
		ja := NewArray([]string{"abc", "xyz"}).Add(3.14)
		ja.Add(27).Add(true)
		ja.Add(New().Add("pi", 3.14))
		ja.Add(uint(123))
		ja.Add(NewArray(123))
		if len(ja.Values) != 8 {
			t.Fail()
		}

		s, ok := ja.GetString(0)
		if !ok || s != "abc" {
			t.Fail()
		}
		if _, ok := ja.GetString(123); ok {
			t.Fail()
		}
		if _, ok := ja.GetString(3); ok {
			t.Fail()
		}

		i, ok := ja.GetInt(3)
		if !ok || i != 27 {
			t.Fail()
		}
		if _, ok := ja.GetInt(0); ok {
			t.Fail()
		}
		if _, ok := ja.GetInt(123); ok {
			t.Fail()
		}

		ui, ok := ja.GetUInt(6)
		if !ok || ui != 123 {
			t.Fail()
		}
		if _, ok := ja.GetUInt(0); ok {
			t.Fail()
		}
		if _, ok := ja.GetUInt(123); ok {
			t.Fail()
		}

		f, ok := ja.GetFloat(2)
		if !ok || f != 3.14 {
			t.Fail()
		}
		if _, ok := ja.GetFloat(0); ok {
			t.Fail()
		}
		if _, ok := ja.GetFloat(123); ok {
			t.Fail()
		}

		b, ok := ja.GetBool(4)
		if !ok || !b {
			t.Fail()
		}
		if _, ok := ja.GetBool(0); ok {
			t.Fail()
		}
		if _, ok := ja.GetBool(123); ok {
			t.Fail()
		}

		jo, ok := ja.GetObject(5)
		if !ok || len(jo.Properties) != 1 {
			t.Fail()
		}
		if _, ok := ja.GetObject(0); ok {
			t.Fail()
		}
		if _, ok := ja.GetObject(123); ok {
			t.Fail()
		}

		if _, ok := ja.GetArray(7); !ok {
			t.Fail()
		}
		if _, ok := ja.GetArray(0); ok {
			t.Fail()
		}
		if _, ok := ja.GetArray(123); ok {
			t.Fail()
		}

		ja.Remove(4)
		if len(ja.Values) != 7 {
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		var ja *Array
		s := ja.String()
		if s != "" {
			t.Fail()
		}

		ja = new(Array)
		s = ja.String()
		if s != "[]" {
			t.Fail()
		}

	})
}
