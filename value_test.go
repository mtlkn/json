package json

import (
	"fmt"
	"testing"
)

func TestValue(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		v := New("YM")
		if v.Type() != STRING || v.Value() != "YM" {
			t.Fail()
			return
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
