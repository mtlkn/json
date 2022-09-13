package json

import (
	"errors"
	"io"
	"strings"
)

type Array struct {
	Values []*Value
}

func NewArray(values ...interface{}) *Array {
	ja := new(Array)

	if len(values) > 0 {
		ja.Values = make([]*Value, len(values))
		for i, v := range values {
			ja.Values[i] = New(v)
		}
	}

	return ja
}

// shortcut for NewArray
func A(values ...interface{}) *Array {
	return NewArray(values...)
}

func (ja *Array) GetStrings() ([]string, bool) {
	if len(ja.Values) == 0 {
		return nil, true
	}

	vs := make([]string, len(ja.Values))

	for i, v := range ja.Values {
		s, ok := v.String()
		if !ok {
			return nil, false
		}
		vs[i] = s
	}

	return vs, true
}

func (ja *Array) GetInts() ([]int, bool) {
	if len(ja.Values) == 0 {
		return nil, true
	}

	vs := make([]int, len(ja.Values))

	for i, v := range ja.Values {
		n, ok := v.Int()
		if !ok {
			return nil, false
		}
		vs[i] = n
	}

	return vs, true
}

func (ja *Array) GetFloats() ([]float64, bool) {
	if len(ja.Values) == 0 {
		return nil, true
	}

	vs := make([]float64, len(ja.Values))

	for i, v := range ja.Values {
		f, ok := v.Float()
		if !ok {
			return nil, false
		}
		vs[i] = f
	}

	return vs, true
}

func (ja *Array) GetObjects() ([]*Object, bool) {
	if len(ja.Values) == 0 {
		return nil, true
	}

	vs := make([]*Object, len(ja.Values))

	for i, v := range ja.Values {
		o, ok := v.Object()
		if !ok {
			return nil, false
		}
		vs[i] = o
	}

	return vs, true
}

func (ja *Array) String() string {
	var sb strings.Builder

	sb.WriteByte('[')

	for i, v := range ja.Values {
		if i > 0 {
			sb.WriteByte(',')
		}

		if v == nil {
			sb.WriteString("null")
			continue
		}

		if len(v.buf) > 0 {
			sb.Write(v.buf)
			continue
		}

		sb.WriteString(v.string())
	}

	sb.WriteByte(']')

	return sb.String()
}

func (ja *Array) Validate() error {
	if ja == nil {
		return nil
	}

	for _, v := range ja.Values {
		err := v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseArray(r io.Reader) (*Array, error) {
	rd, err := newReader(r)
	if err != nil {
		return nil, err
	}

	if !rd.SkipSpace() || rd.b != '[' {
		err = rd.EnsureJSON('[')
		if err != nil {
			return nil, err
		}
	}

	return rd.parseArray()
}

func (rd *reader) parseArray() (*Array, error) {
	ja := new(Array)

	for {
		if !rd.SkipSpace() {
			return nil, errors.New("missing closing ]")
		}

		if rd.b == ']' {
			return ja, nil
		}

		var v *Value

		skip := true

		switch rd.b {
		case '"':
			l := rd.i
			var ok bool
			for rd.Read() {
				if rd.b == '"' && rd.buf[rd.i-1] != '\\' {
					v = &Value{
						buf: rd.buf[l : rd.i+1],
					}
					ok = true
					break
				}
			}
			if !ok {
				return nil, errors.New("missing closing quote")
			}
		case '{':
			jo, err := rd.parseObject()
			if err != nil {
				return nil, err
			}
			v = &Value{
				typ: OBJECT,
				val: jo,
			}
		case '[':
			ja, err := rd.parseArray()
			if err != nil {
				return nil, err
			}
			v = &Value{
				typ: ARRAY,
				val: ja,
			}
		default:
			l := rd.i
			var ok bool
			for rd.Read() {
				skip = IsSpace(rd.b)
				if rd.b == ',' || rd.b == ']' || skip {
					v = &Value{
						buf: rd.buf[l:rd.i],
					}
					ok = true
					break
				}
			}
			if !ok {
				return nil, errors.New("missing closing ]")
			}
		}

		if skip && !rd.SkipSpace() {
			return nil, errors.New("missing closing ]")
		}

		ja.Values = append(ja.Values, v)

		if rd.b == ']' {
			break
		}

		if rd.b != ',' {
			return nil, errors.New("missing array value closing")
		}
	}

	return ja, nil
}
