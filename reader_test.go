package json

import (
	"errors"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	s := ` { "name": "Yuri Metelkin \"YM\"" } `
	r, err := newReader(strings.NewReader(s))
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}

	if len(r.buf) != len(s) {
		t.Fail()
		return
	}

	if !r.SkipSpaceTo('{') || r.b != '{' {
		t.Fail()
		return
	}

	if !r.SkipSpaceTo('"') || r.b != '"' {
		t.Fail()
		return
	}

	q, ok := r.ReadQuotes()
	if !ok || q != "name" {
		t.Fail()
		return
	}

	if !r.SkipSpaceTo(':') || r.b != ':' {
		t.Fail()
		return
	}

	if !r.SkipSpaceTo('"') || r.b != '"' {
		t.Fail()
		return
	}

	q, ok = r.ReadQuotes()
	if !ok || q != "Yuri Metelkin \\\"YM\\\"" {
		t.Fail()
		return
	}

	if !r.SkipSpaceTo('}') || r.b != '}' {
		t.Fail()
		return
	}

	r, _ = newReader(strings.NewReader("\""))
	r.SkipSpace()
	if r.b != '"' {
		t.Fail()
		return
	}
	_, ok = r.ReadQuotes()
	if ok {
		t.Fail()
		return
	}

	ok = r.SkipSpaceTo('"')
	if ok {
		t.Fail()
	}

	_, err = newReader(new(errorReader))
	if err == nil {
		t.Fail()
	}
}

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("dummy error")
}
