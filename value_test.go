package json

import "testing"

func TestValue(t *testing.T) {
	t.Run("new value", func(t *testing.T) {
		v := newValue("xyz")
		s, ok := v.GetString()
		if !ok || s != "xyz" || v.Type != STRING {
			t.Fail()
		}

		v = newValue(1)
		d, ok := v.GetInt()
		if !ok || d != 1 || v.Type != INT {
			t.Fail()
		}

		if _, ok = v.GetUInt(); ok {
			t.Fail()
		}
		v = newValue(uint(1))
		ui, ok := v.GetUInt()
		if !ok || ui != 1 || v.Type != UINT {
			t.Fail()
		}

		v = newValue(3.14)
		f, ok := v.GetFloat()
		if !ok || f != 3.14 || v.Type != FLOAT {
			t.Fail()
		}

		v = newValue(true)
		b, ok := v.GetBool()
		if !ok || !b || v.Type != BOOL {
			t.Fail()
		}

		v = newValue([]string{"abc", "xyz"})
		if v.Type != ARRAY || v.String() != `["abc","xyz"]` {
			t.Fail()
		}

		v = newValue([]int{1, 2})
		if v.Type != ARRAY || v.String() != "[1,2]" {
			t.Fail()
		}

		v = newValue([]uint{1, 2})
		if v.Type != ARRAY || v.String() != "[1,2]" {
			t.Fail()
		}

		v = newValue([]float64{3.14, 0.2e-3})
		if v.Type != ARRAY || v.String() != "[3.14,0.0002]" {
			t.Fail()
		}

		v = newValue([]bool{true, false})
		if v.Type != ARRAY || v.String() != "[true,false]" {
			t.Fail()
		}

		v = newValue([]*Object{New().Add("name", "YM")})
		if v.Type != ARRAY || v.String() != `[{"name":"YM"}]` {
			t.Fail()
		}

		v = newValue(new(Property))
		if v.Type != NULL || v.String() != "null" {
			t.Fail()
		}

		v = &Value{
			Type:  NULL,
			value: 123,
		}
		if v.String() != "null" {
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		v := &Value{
			Type:  OBJECT,
			value: 1,
		}
		if v.String() != "" {
			t.Fail()
		}

		v = &Value{
			Type:  ARRAY,
			value: 1,
		}
		if v.String() != "" {
			t.Fail()
		}

		v = &Value{}
		if _, err := v.GetValue(); err == nil {
			t.Fail()
		}

		v = &Value{
			data: []byte("abc"),
		}
		if _, err := v.GetValue(); err == nil {
			t.Fail()
		}

		v = &Value{
			data:    []byte("12.34.56.87"),
			special: floatBytes,
		}
		if _, err := v.GetValue(); err == nil {
			t.Fail()
		}
	})
}
