package json

import (
	"bytes"
	"os"
	"testing"
)

func BenchmarkPrimeNumbers(b *testing.B) {
	fs, err := os.ReadDir("testdata")
	if err != nil {
		b.Error("failed to read testdata")
		b.Fail()
		return
	}

	for _, f := range fs {
		b.Run(f.Name(), func(b *testing.B) {
			bs, _ := os.ReadFile("testdata/" + f.Name())

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := ParseObject(bytes.NewReader(bs))
				if err != nil {
					b.Error(err)
					b.Fail()
					return
				}
			}
		})
	}
}
