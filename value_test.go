package json

import (
	"fmt"
	"testing"
)

type KVP map[string]interface{}

func (kvp KVP) O() *Object {
	o := O()
	for k, v := range kvp {
		o.Set(k, v)
	}
	return o
}

func TestXxx(t *testing.T) {
	o := KVP{
		"id":   1,
		"name": "YM",
	}.O()
	fmt.Println(o.String())
}

func TestValue(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		v := New("YM")
		if v.Type() != STRING || v.Value() != "YM" {
			t.Fail()
			return
		}
	})

	t.Run("strings value", func(t *testing.T) {
		v := New([]string{"YM", "SV"})
		if v.Type() != ARRAY {
			t.Fail()
			return
		}

		ja, ok := v.Array()
		if !ok {
			t.Fail()
			return
		}

		ss, ok := ja.GetStrings()
		if !ok || len(ss) != 2 || ss[0] != "YM" || ss[1] != "SV" {
			t.Fail()
		}
	})

	t.Run("int value", func(t *testing.T) {
		i := 1
		for _, value := range []interface{}{i, int8(i), int16(i), int32(i), int64(i)} {
			v := New(value)
			if v.Type() != INT || v.Value() != i {
				t.Fail()
				return
			}

			f, ok := v.Float()
			if !ok || f != 1 {
				t.Fail()
				return
			}

			u, ok := v.UInt()
			if !ok || u != 1 {
				t.Fail()
				return
			}
		}
	})

	t.Run("ints value", func(t *testing.T) {
		v := New([]int{1, 2})
		if v.Type() != ARRAY {
			t.Fail()
			return
		}

		ja, ok := v.Array()
		if !ok {
			t.Fail()
			return
		}

		is, ok := ja.GetInts()
		if !ok || len(is) != 2 || is[0] != 1 || is[1] != 2 {
			t.Fail()
		}

		for _, value := range []interface{}{[]int8{1}, []int16{1}, []int32{1}, []int64{1}} {
			v := New(value)
			if v.Type() != ARRAY {
				t.Fail()
				return
			}

			ja, ok := v.Array()
			if !ok {
				t.Fail()
				return
			}

			is, ok := ja.GetInts()
			if !ok || len(is) != 1 || is[0] != 1 {
				t.Fail()
			}
		}
	})

	t.Run("float value", func(t *testing.T) {
		f := 3.14
		for _, value := range []interface{}{f, float32(f)} {
			v := New(value)
			if v.Type() != FLOAT || v.Value() != f {
				t.Fail()
				return
			}

			i, ok := v.Int()
			if !ok || i != 3 {
				t.Fail()
				return
			}

			u, ok := v.UInt()
			if !ok || u != 3 {
				t.Fail()
				return
			}
		}
	})

	t.Run("floats value", func(t *testing.T) {
		for _, value := range []interface{}{[]float64{.01, 3.14}, []float32{.01, 3.14}} {
			v := New(value)

			if v.Type() != ARRAY {
				t.Fail()
				return
			}

			ja, ok := v.Array()
			if !ok {
				t.Fail()
				return
			}

			fs, ok := ja.GetFloats()
			if !ok || len(fs) != 2 || fs[0] != 0.01 || fs[1] != 3.14 {
				t.Fail()
			}
		}
	})

	t.Run("uint value", func(t *testing.T) {
		var i uint = 1
		for _, value := range []interface{}{i, uint8(i), uint16(i), uint32(i), uint64(i)} {
			v := New(value)
			if v.Type() != UINT || v.Value() != i {
				t.Fail()
				return
			}
		}
	})

	t.Run("uints value", func(t *testing.T) {
		for _, value := range []interface{}{[]uint{1}, []uint8{1}, []uint16{1}, []uint32{1}, []uint64{1}} {
			v := New(value)
			if v.Type() != ARRAY {
				t.Fail()
				return
			}

			ja, ok := v.Array()
			if !ok || len(ja.Values) != 1 || ja.Values[0].Value() != uint(1) {
				t.Fail()
			}
		}
	})

	t.Run("bool value", func(t *testing.T) {
		v := New(true)
		if v.Type() != BOOL || v.Value() != true {
			t.Fail()
			return
		}
	})

	t.Run("object value", func(t *testing.T) {
		v := New(O(P("name", "YM")))
		if v.Type() != OBJECT {
			t.Fail()
			return
		}

		o, ok := v.Object()
		if !ok {
			t.Fail()
			return
		}

		s, ok := o.GetString("name")
		if !ok || s != "YM" {
			t.Fail()
			return
		}
	})

	t.Run("objects value", func(t *testing.T) {
		v := New([]*Object{
			O(P("id", 1), P("name", "YM")),
			O(P("id", 2), P("name", "SV")),
		})
		if v.Type() != ARRAY {
			t.Fail()
			return
		}

		ja, ok := v.Array()
		if !ok {
			t.Fail()
			return
		}

		oo, ok := ja.GetObjects()
		if !ok || len(oo) != 2 {
			t.Fail()
		}

		if id, _ := oo[0].GetInt("id"); id != 1 {
			t.Fail()
		}

		if s, _ := oo[1].GetString("name"); s != "SV" {
			t.Fail()
		}
	})

	t.Run("array value", func(t *testing.T) {
		v := New(A("YM", "SV"))
		if v.Type() != ARRAY {
			t.Fail()
			return
		}

		a, ok := v.Array()
		if !ok || len(a.Values) != 2 {
			t.Fail()
			return
		}

		ss, ok := a.GetStrings()
		if !ok || len(ss) != 2 || ss[0] != "YM" || ss[1] != "SV" {
			t.Fail()
			return
		}
	})

	t.Run("value string", func(t *testing.T) {
		var v *Value
		if v.string() != "null" {
			t.Fail()
		}

		v = new(Value)
		if v.string() != "null" {
			t.Fail()
		}

		v.typ = OBJECT
		if v.string() != "null" {
			t.Fail()
		}

		v.typ = ARRAY
		if v.string() != "null" {
			t.Fail()
		}

		v.typ = BOOL
		if v.string() != "false" {
			t.Fail()
		}

		v = &Value{
			typ: UINT,
			val: uint(1),
		}
		if v.string() != "1" {
			t.Fail()
		}
	})
}

func TestValueParse(t *testing.T) {
	t.Run("parsing string value", func(t *testing.T) {
		s := `"text"`

		v := &Value{
			buf: []byte(s),
		}
		if v.Type() != STRING {
			t.Errorf("wrong type: expected %v, got %v", STRING, v.Type())
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		if v.Value() != "text" {
			t.Errorf("wrong value: expected %v, got %v", "text", v.Value())
			t.Fail()
		}

		s, ok := v.String()
		if !ok || s != "text" {
			t.Errorf("wrong value: expected %v, got %v", "text", v.Value())
			t.Fail()
		}

		if v.Debug() != `{"type":3,"parsed":""text"","value":text}` {
			t.Fail()
		}

		v = &Value{
			buf: []byte("\""),
		}
		ok = v.parseString()
		if ok || v.Type() != INVALID {
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		_, ok = v.UInt()
		if ok {
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		_, ok = v.UInt()
		if ok {
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		_, ok = v.Int()
		if ok {
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		_, ok = v.Float()
		if ok {
			t.Fail()
		}
		_, ok = v.Bool()
		if ok {
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		_, ok = v.Object()
		if ok {
			t.Fail()
		}

		v = &Value{
			buf: []byte(s),
		}
		_, ok = v.Array()
		if ok {
			t.Fail()
		}
	})

	t.Run("parsing int value", func(t *testing.T) {
		s := "1965"
		v := &Value{
			buf: []byte(s),
		}
		err := v.Validate()
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}
		if v.Type() != INT {
			t.Errorf("wrong type: expected %v, got %v", INT, v.Type())
			t.Fail()
		}
		if v.Value() != 1965 {
			t.Errorf("wrong value: expected %v, got %v", 1965, v.Value())
			t.Fail()
		}
		i, ok := v.Int()
		if !ok || i != 1965 {
			t.Errorf("wrong int value: expected %v, got %v", 1965, i)
			t.Fail()
		}

		s = "-25"
		v = &Value{
			buf: []byte(s),
		}
		err = v.Validate()
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}
		if v.Type() != INT {
			t.Errorf("wrong type: expected %v, got %v", INT, v.Type())
			t.Fail()
		}
		if v.Value() != -25 {
			t.Errorf("wrong value: expected %v, got %v", -25, v.Value())
			t.Fail()
		}
		i, ok = v.Int()
		if !ok || i != -25 {
			t.Errorf("wrong int value: expected %v, got %v", -25, i)
			t.Fail()
		}

		s = "+30"
		v = &Value{
			buf: []byte(s),
		}
		err = v.Validate()
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}
		if v.Type() != INT {
			t.Errorf("wrong type: expected %v, got %v", INT, v.Type())
			t.Fail()
		}
		if v.Value() != 30 {
			t.Errorf("wrong value: expected %v, got %v", 30, v.Value())
			t.Fail()
		}
		i, ok = v.Int()
		if !ok || i != 30 {
			t.Errorf("wrong int value: expected %v, got %v", 30, i)
			t.Fail()
		}

		s = "18446744073709551615"
		ui := uint(18446744073709551615)
		v = &Value{
			buf: []byte(s),
		}
		if v.Type() != UINT {
			t.Errorf("wrong type: expected %v, got %v", UINT, v.Type())
			t.Fail()
		}
		if v.Value() != ui {
			t.Errorf("wrong value: expected %v, got %v", ui, v.Value())
			t.Fail()
		}
		u, ok := v.UInt()
		if !ok || u != ui {
			t.Errorf("wrong int value: expected %v, got %v", ui, u)
			t.Fail()
		}

		s = "6666666666666666666666666666666666666666666666666666666666666666666"
		v = &Value{
			buf: []byte(s),
		}
		if v.Type() != INVALID {
			t.Fail()
		}
	})

	t.Run("parsing float value", func(t *testing.T) {
		for _, s := range []string{"3.14", "1.2345E+13", "1.2345e-13", "-3.14"} {
			v := &Value{
				buf: []byte(s),
			}
			err := v.Validate()
			if err != nil {
				t.Fail()
				t.Error(err)
				return
			}

			if v.Type() != FLOAT {
				t.Errorf("wrong type: expected %v, got %v", FLOAT, v.Type())
				t.Fail()
			}

			f, ok := v.Float()
			if !ok {
				t.Fail()
			}
			fmt.Println(f)
		}

		v := &Value{
			buf: []byte("1.2.3.4"),
		}
		err := v.Validate()
		if err == nil {
			t.Fail()
			return
		}
	})

	t.Run("parsing false value", func(t *testing.T) {
		s := "false"
		v := &Value{
			buf: []byte(s),
		}
		err := v.Validate()
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}

		if v.Type() != BOOL {
			t.Errorf("wrong type: expected %v, got %v", BOOL, v.Type())
			t.Fail()
		}
		if v.Value() == true {
			t.Errorf("wrong value: expected %v, got %v", false, v.Value())
			t.Fail()
		}
		b, ok := v.Bool()
		if !ok || b {
			t.Errorf("wrong int value: expected %v, got %v", false, b)
			t.Fail()
		}
	})

	t.Run("parsing null value", func(t *testing.T) {
		s := "null"
		v := &Value{
			buf: []byte(s),
		}
		err := v.Validate()
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}

		if v.Type() != NULL {
			t.Errorf("wrong type: expected %v, got %v", NULL, v.Type())
			t.Fail()
		}
		if v.Value() != nil {
			t.Errorf("wrong value: expected nil, got %v", v.Value())
			t.Fail()
		}
	})
}
