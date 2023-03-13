package json

import (
	"errors"
	"fmt"
)

func Parse(bs []byte) (*Value, error) {
	l, r, t, err := trimJSONInput(bs)
	if err != nil {
		return nil, err
	}

	if t == OBJECT {
		jo, _, err := parseObject(bs, uint(l+1), uint(r+1))
		return &Value{
			Type:  OBJECT,
			value: jo,
		}, err
	}

	ja, _, err := parseArray(bs, uint(l+1), uint(r+1))
	return &Value{
		Type:  ARRAY,
		value: ja,
	}, err
}

func ParseString(s string) (*Value, error) {
	return Parse(stringToBytes(s))
}

func ParseObject(bs []byte) (*Object, error) {
	l, r, t, err := trimJSONInput(bs)
	if err != nil {
		return nil, err
	}
	if t != OBJECT {
		return nil, errors.New("bad json object")
	}
	jo, _, err := parseObject(bs, uint(l+1), uint(r+1))
	return jo, err
}

func ParseObjectString(s string) (*Object, error) {
	return ParseObject(stringToBytes(s))
}

func ParseArray(bs []byte) (*Array, error) {
	l, r, t, err := trimJSONInput(bs)
	if err != nil {
		return nil, err
	}
	if t != ARRAY {
		return nil, errors.New("bad json array")
	}
	ja, _, err := parseArray(bs, uint(l+1), uint(r+1))
	return ja, err
}

func ParseArrayString(s string) (*Array, error) {
	return ParseArray(stringToBytes(s))
}

func parseObject(bs []byte, start, last uint) (*Object, uint, error) {
	if start >= last {
		return nil, start, fmt.Errorf("object parsing at %d", start)
	}

	var jps []*Property

	i := start
	for i < last {
		jp, end, err := parseProperty(bs, i, last)
		if err != nil {
			return nil, end, err
		}

		jps = append(jps, jp)

		b := bs[end]

		if b == ',' {
			i = end + 1
			continue
		}

		if b == '}' {
			return &Object{
				Properites: jps,
			}, end, nil
		}

		break
	}

	return nil, i, fmt.Errorf("object parsing at %d", i)
}

func parseProperty(bs []byte, start, last uint) (*Property, uint, error) {
	if start >= last {
		return nil, start, fmt.Errorf("property parsing at %d", start)
	}

	// parse name
	var name string

	for i := start; i < last; i++ {
		if bs[i] == '"' {
			end, err := parseString(bs, i, last)
			if err != nil {
				return nil, end, err
			}

			name = bytesToString(bs[i+1 : end])
			start = end + 1
			break
		} else if !isWS(bs[i]) {
			return nil, i, fmt.Errorf("property name parsing at %d", i)
		}
	}

	if start >= last {
		return nil, start, fmt.Errorf("property name parsing at %d", start)
	}

	// parse colon
	for i := start; i < last; i++ {
		if bs[i] == ':' {
			start = i + 1
			break
		} else if !isWS(bs[i]) {
			return nil, i, fmt.Errorf("property parsing at %d", i)
		}
	}

	if start >= last {
		return nil, start, fmt.Errorf("property parsing at %d", start)
	}

	value, end, err := parseValue(bs, start, last)
	if err != nil {
		return nil, end, err
	}

	return &Property{
		Name:  name,
		Value: value,
	}, end, nil
}

func parseArray(bs []byte, start, last uint) (*Array, uint, error) {
	if start >= last {
		return nil, start, fmt.Errorf("array parsing at %d", start)
	}

	var vs []*Value

	i := start
	for i < last {
		v, end, err := parseValue(bs, i, last)
		if err != nil {
			return nil, end, err
		}

		vs = append(vs, v)

		b := bs[end]

		if b == ',' {
			i = end + 1
			continue
		}

		if b == ']' {
			return &Array{
				Values: vs,
			}, end, nil
		}

		break
	}

	return nil, i, fmt.Errorf("array parsing at %d", i)
}

func parseValue(bs []byte, start, last uint) (*Value, uint, error) {
	if start >= last {
		return nil, start, fmt.Errorf("value parsing at %d", start)
	}

	var value *Value

	for i := start; i < last; i++ {
		b := bs[i]

		if isWS(b) {
			continue
		}

		if b == '"' {
			end, special, err := parseStringValue(bs, i, last)
			if err != nil {
				return nil, end, err
			}

			value = &Value{
				Type:    STRING,
				data:    bs[i : end+1],
				special: special,
			}

			start = end + 1
			break
		}

		if b == '{' {
			jo, end, err := parseObject(bs, i+1, last)
			if err != nil {
				return nil, end, err
			}

			value = &Value{
				Type:  OBJECT,
				value: jo,
			}

			start = end + 1
			break
		}

		if b == '[' {
			ja, end, err := parseArray(bs, i+1, last)
			if err != nil {
				return nil, end, err
			}

			value = &Value{
				Type:  ARRAY,
				value: ja,
			}

			start = end + 1
			break
		}

		start = i + 1
		if start >= last {
			return nil, start, fmt.Errorf("value parsing at %d", start)
		}

		end, special, err := parseNonTextValue(bs, start, last)
		if err != nil {
			return nil, start, err
		}

		value = &Value{
			data:    bs[i:end],
			special: special,
		}
		start = end
		break
	}

	// parse separator
	for i := start; i < last; i++ {
		if !isWS(bs[i]) {
			return value, i, nil
		}
	}

	return nil, start, fmt.Errorf("value parsing at %d", start)
}

func parseString(bs []byte, start, last uint) (uint, error) {
	if start >= last {
		return start, fmt.Errorf("string parsing at %d", start)
	}

	for i := start + 1; i < last; i++ {
		if bs[i] == '"' && bs[i-1] != '\\' {
			return i, nil
		}
	}

	return start, fmt.Errorf("string parsing at %d", start)
}

func parseStringValue(bs []byte, start, last uint) (uint, specialBytes, error) {
	if start >= last {
		return start, 0, fmt.Errorf("string parsing at %d", start)
	}

	var special specialBytes

	for i := start + 1; i < last; i++ {
		switch bs[i] {
		case '"':
			if bs[i-1] != '\\' {
				return i, special, nil
			}
		case '\\':
			special = unquoteBytes
		}
	}

	return start, 0, fmt.Errorf("string parsing at %d", start)
}

func parseNonTextValue(bs []byte, start, last uint) (uint, specialBytes, error) {
	var special specialBytes

	for i := start; i < last; i++ {
		b := bs[i]

		if b == '.' || b == 'e' || b == 'E' {
			special = floatBytes
		}

		if isWS(b) || b == ',' || b == '}' || b == ']' {
			return i, special, nil
		}
	}

	return start, 0, fmt.Errorf("property value parsing at %d", start)
}

func isWS(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\v' || b == '\f' || b == '\r'
}

func trimJSONInput(bs []byte) (int, int, ValueType, error) {
	last := len(bs) - 1

	if last < 1 {
		return 0, 0, 0, errors.New("bad json")
	}

	for l := 0; l < last; l++ {
		if bs[l] == '{' {
			for r := last; r > l; r-- {
				if bs[r] == '}' {
					return l, r, OBJECT, nil
				}

				if isWS(bs[r]) {
					continue
				}

				return 0, 0, 0, errors.New("bad json")
			}
		}

		if bs[l] == '[' {
			for r := last; r > l; r-- {
				if bs[r] == ']' {
					return l, r, ARRAY, nil
				}

				if isWS(bs[r]) {
					continue
				}

				return 0, 0, 0, errors.New("bad json")
			}
		}

		if isWS(bs[l]) {
			continue
		}

		break
	}

	return 0, 0, 0, errors.New("bad json")
}
