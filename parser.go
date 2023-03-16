package json

import (
	"errors"
	"fmt"
)

var (
	errorInvalidJSON   = errors.New("invalid JSON")
	errorInvalidObject = errors.New("invalid JSON object")
	errorInvalidArray  = errors.New("invalid JSON array")
)

func Parse(bs []byte) (*Value, error) {
	l, r, t, err := trimJSONInput(bs)
	if err != nil {
		return nil, err
	}

	if t == OBJECT {
		jo, _, err := parseObject(bs, l+1, r+1)
		return &Value{
			Type:  OBJECT,
			value: jo,
		}, err
	}

	ja, _, err := parseArray(bs, l+1, r+1)
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
		return nil, errorInvalidObject
	}
	jo, _, err := parseObject(bs, l+1, r+1)
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
		return nil, errorInvalidArray
	}
	ja, _, err := parseArray(bs, l+1, r+1)
	return ja, err
}

func ParseArrayString(s string) (*Array, error) {
	return ParseArray(stringToBytes(s))
}

func parseObject(bs []byte, start, last int) (*Object, int, error) {
	if start >= last {
		return nil, start, parsingError("object", start)
	}

	var jps []*Property

	i := start
	for i < last {
		jp, end, err := parseProperty(bs, i, last)
		if err != nil {
			return nil, end, err
		}

		if jp != nil {
			jps = append(jps, jp)
		}

		b := bs[end]

		if b == ',' {
			i = end + 1
			continue
		}

		if b == '}' {
			return &Object{
				Properties: jps,
			}, end, nil
		}

		break
	}

	return nil, i, parsingError("object", i)
}

func parseProperty(bs []byte, start, last int) (*Property, int, error) {
	if start >= last {
		return nil, start, parsingError("property", start)
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
		} else if bs[i] == '}' {
			return nil, i, nil
		} else if !isWS(bs[i]) {
			return nil, i, parsingError("property name", i)
		}
	}

	if start >= last {
		return nil, start, parsingError("property name", start)
	}

	// parse colon
	for i := start; i < last; i++ {
		if bs[i] == ':' {
			start = i + 1
			break
		} else if !isWS(bs[i]) {
			return nil, i, parsingError("property", i)
		}
	}

	if start >= last {
		return nil, start, parsingError("property", start)
	}

	value, end, err := parseValue(bs, start, last)
	if err != nil {
		return nil, end, err
	}

	if value == nil {
		return nil, end, parsingError("property", start)
	}

	return &Property{
		Name:  name,
		Value: value,
	}, end, nil
}

func parseArray(bs []byte, start, last int) (*Array, int, error) {
	if start >= last {
		return nil, start, parsingError("array", start)
	}

	var vs []*Value

	i := start
	for i < last {
		v, end, err := parseValue(bs, i, last)
		if err != nil {
			return nil, end, err
		}

		if v != nil {
			vs = append(vs, v)
		}

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

	return nil, i, parsingError("array", i)
}

func parseValue(bs []byte, start, last int) (*Value, int, error) {
	if start >= last {
		return nil, start, parsingError("value", start)
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

		if b == '}' || b == ']' {
			start = i
			break
		}

		start = i + 1
		if start >= last {
			return nil, start, parsingError("value", i)
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

	return nil, start, parsingError("value", start)
}

func parseString(bs []byte, start, last int) (int, error) {
	if start >= last {
		return start, parsingError("string", start)
	}

	for i := start + 1; i < last; i++ {
		if bs[i] == '"' && bs[i-1] != '\\' {
			return i, nil
		}
	}

	return start, parsingError("string", start)
}

func parseStringValue(bs []byte, start, last int) (int, specialBytes, error) {
	if start >= last {
		return start, 0, parsingError("string value", start)
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

	return start, 0, parsingError("string value", start)
}

func parseNonTextValue(bs []byte, start, last int) (int, specialBytes, error) {
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

	return start, 0, parsingError("property value", start)
}

func isWS(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\v' || b == '\f' || b == '\r'
}

func trimJSONInput(bs []byte) (int, int, ValueType, error) {
	last := len(bs) - 1

	if last < 1 {
		return 0, 0, 0, errorInvalidJSON
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

				return 0, 0, 0, errorInvalidJSON
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

				return 0, 0, 0, errorInvalidJSON
			}
		}

		if isWS(bs[l]) {
			continue
		}

		break
	}

	return 0, 0, 0, errorInvalidJSON
}

func parsingError(what string, where int) error {
	return fmt.Errorf("%s parsing at %d", what, where)
}
