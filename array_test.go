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
