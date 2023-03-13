package json

import (
	"testing"
)

func BenchmarkBytes(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	bs := []byte(s)

	b.Run("safe bytes to string", func(b *testing.B) {
		fn := func() string {
			return string(bs)
		}
		for n := 0; n < b.N; n++ {
			fn()
		}
	})

	b.Run("safe string to bytes", func(b *testing.B) {
		fn := func() []byte {
			return []byte(s)
		}
		for n := 0; n < b.N; n++ {
			fn()
		}
	})

	b.Run("unsafe bytes to string", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			bytesToString(bs)
		}
	})

	b.Run("unsafe string to bytes", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			stringToBytes(s)
		}
	})
}

func TestBytes(t *testing.T) {
	t.Run("string to bytes", func(t *testing.T) {
		s := "The quick brown fox jumps over the lazy dog"
		bs := stringToBytes(s)

		if len(bs) != len(s) {
			t.Fail()
		}

		if string(bs) != s {
			t.Fail()
		}

		bs = stringToBytes("")
		if len(bs) > 0 {
			t.Fail()
		}
	})

	t.Run("bytes to string", func(t *testing.T) {
		bs := []byte("The quick brown fox jumps over the lazy dog")
		s := bytesToString(bs)

		if len(bs) != len(s) {
			t.Fail()
		}
		if string(bs) != s {
			t.Fail()
		}

		s = bytesToString(nil)
		if s != "" {
			t.Fail()
		}
	})
}
