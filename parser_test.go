package json

import (
	"testing"
)

func TestParser(t *testing.T) {
	p := newParser([]byte{})
	if err := p.EnsureJSON(); err != nil {
		t.Fail()
		t.Error("failed to ensure empty json")
	}

	p = newParser([]byte{'{'})
	if err := p.EnsureJSON(); err == nil {
		t.Fail()
		t.Error("failed to fail to ensure bad json")
	}
}

var (
	data = make([][]byte, 3)
	keys = []string{"hits", "messi", "search"}
)
