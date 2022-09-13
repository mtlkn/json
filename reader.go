package json

import (
	"io"
)

type reader struct {
	b   byte   // current byte
	i   int    // current index
	buf []byte // byte buffer
	sz  int    // buffer size
}

func newReader(r io.Reader) (*reader, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	sz := len(buf)
	if sz == 0 {
		return nil, io.EOF
	}

	return &reader{
		i:   -1,
		buf: buf,
		sz:  len(buf),
	}, nil
}

func (rd *reader) Read() bool {
	rd.i++
	if rd.i >= rd.sz {
		return false
	}

	rd.b = rd.buf[rd.i]
	return true
}

func (rd *reader) SkipSpace() bool {
	for rd.Read() {
		if IsSpace(rd.b) {
			continue
		}
		return true
	}

	return false
}

func (rd *reader) SkipSpaceTo(b byte) bool {
	if !rd.SkipSpace() {
		return false
	}
	return rd.b == b
}

/*
func (rd *reader) ReadTo(b byte, gte, lte bool) ([]byte, bool) {
	l := rd.i
	if !gte {
		l++
	}

	for rd.Read() {
		if rd.b == b {
			r := rd.i
			if lte {
				r++
			}
			return rd.buf[l:r], true
		}
	}

	return nil, false
}
*/

func (rd *reader) ReadQuotes() (string, bool) {
	l := rd.i + 1

	for rd.Read() {
		if rd.b == '"' && rd.buf[rd.i-1] != '\\' {
			buf := rd.buf[l:rd.i]
			return BytesToString(buf), true
		}
	}

	return "", false
}

func (rd *reader) EnsureJSON(opening byte) error {
	if rd.sz < 2 {
		return ErrBadJSON
	}

	var last byte

	rd.i = -1

	for rd.Read() {
		if rd.b == '{' {
			if opening != rd.b {
				return ErrBadJSON
			}
			last = '}'
			break
		}

		if rd.b == '[' {
			if opening != rd.b {
				return ErrBadJSON
			}
			last = ']'
			break
		}

		if rd.b > 32 && rd.b < 127 {
			return ErrBadJSON
		}
	}

	if last == 0 {
		return ErrBadJSON
	}

	l := rd.i // opening JSON

	r := rd.sz - 1
	for r > 0 {
		b := rd.buf[r]

		if b == last {
			rd.buf = rd.buf[l : r+1]
			rd.sz = len(rd.buf)
			rd.i = -1
			rd.b = 0
			rd.Read()
			return nil
		}

		if b > 32 && b < 127 {
			break
		}

		r--
	}

	return ErrBadJSON
}
