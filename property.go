package json

import (
	"errors"
)

type Property struct {
	n    []byte // parsed name bytes
	name string
	v    []byte // parsed value bytes
	val  *Value
}

func NewProperty(name string, value interface{}) *Property {
	return &Property{
		name: name,
		val:  New(value),
	}
}

// shortcut for NewProperty
func P(name string, value interface{}) *Property {
	return NewProperty(name, value)
}

func (p *Property) Name() string {
	if p.name == "" && len(p.n) > 1 {
		p.name = BytesToString(p.n[1 : len(p.n)-1])
	}
	return p.name
}

func (p *Property) Value() *Value {
	if p.val == nil {
		p.val = &Value{
			buf: p.v,
		}
	}
	return p.val
}

func (rd *reader) parseProperty() (*Property, error) {
	if !rd.SkipSpace() {
		return nil, errors.New("missing closing }")
	}

	if rd.b == '}' {
		return nil, nil
	}

	if rd.b != '"' {
		return nil, errors.New("missing property openning quote")
	}

	p := new(Property)
	l := rd.i

	for rd.Read() {
		if rd.b != '"' || rd.buf[rd.i-1] == '\\' {
			continue
		}

		p.n = rd.buf[l : rd.i+1]

		if !rd.SkipSpace() || rd.b != ':' {
			return nil, errors.New("missing property colon punctuation")
		}

		if !rd.SkipSpace() || rd.b == ',' || rd.b == '}' {
			return nil, errors.New("missing property value")
		}

		skip := true

		switch rd.b {
		case '"':
			l = rd.i
			for rd.Read() {
				if rd.b == '"' && rd.buf[rd.i-1] != '\\' {
					p.v = rd.buf[l : rd.i+1]
					break
				}
			}
		case '{':
			jo, err := rd.parseObject()
			if err != nil {
				return nil, err
			}
			p.val = &Value{
				typ: OBJECT,
				val: jo,
			}
		case '[':
			ja, err := rd.parseArray()
			if err != nil {
				return nil, err
			}
			p.val = &Value{
				typ: ARRAY,
				val: ja,
			}
		default:
			l = rd.i
			for rd.Read() {
				skip = IsSpace(rd.b)
				if rd.b == ',' || rd.b == '}' || skip {
					p.v = rd.buf[l:rd.i]
					break
				}
			}
		}

		if skip && !rd.SkipSpace() {
			return nil, errors.New("missing closing }")
		} else {
			break
		}
	}

	return p, nil
}
