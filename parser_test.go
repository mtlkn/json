package json

import "testing"

func TestParser(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		if _, err := ParseString("{"); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(" a "); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString("{"); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(" [ ] "); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(" [ a "); err == nil {
			t.Fail()
		}

		if _, err := ParseArrayString(" a "); err == nil {
			t.Fail()
		}

		if _, err := ParseArrayString("["); err == nil {
			t.Fail()
		}

		if _, err := ParseArrayString(" { } "); err == nil {
			t.Fail()
		}

		if _, err := ParseArrayString(" { a "); err == nil {
			t.Fail()
		}

		if _, _, err := parseObject([]byte(" {}"), 1, 1); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name:}`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{name:YM}`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name":}`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name": {}`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name": [}`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name": YM`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name "YM"}`); err == nil {
			t.Fail()
		}

		if _, err := ParseObjectString(`{"name": "YM" ?}`); err == nil {
			t.Fail()
		}

		if _, _, err := parseProperty(nil, 1, 1); err == nil {
			t.Fail()
		}

		if _, _, err := parseProperty([]byte(`{"name":"YM"}`), 1, 7); err == nil {
			t.Fail()
		}

		if _, _, err := parseProperty([]byte(`{"name":"YM"}`), 1, 8); err == nil {
			t.Fail()
		}

		if _, _, err := parseArray(nil, 1, 1); err == nil {
			t.Fail()
		}

		if _, err := ParseArrayString(`["abc, xyz]`); err == nil {
			t.Fail()
		}

		if _, err := ParseArrayString(`["abc" xyz]`); err == nil {
			t.Fail()
		}

		if _, _, err := parseValue(nil, 1, 1); err == nil {
			t.Fail()
		}

		if _, err := parseString(nil, 1, 1); err == nil {
			t.Fail()
		}

		if _, _, err := parseStringValue(nil, 1, 1); err == nil {
			t.Fail()
		}

		if _, _, err := parseValue([]byte("123"), 0, 3); err == nil {
			t.Fail()
		}

		if _, _, err := parseValue([]byte("123 "), 0, 4); err == nil {
			t.Fail()
		}

	})
}
