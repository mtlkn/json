package json

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

var (
	ErrBadJSON error = errors.New("invalid JSON")
)

func IsSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r' || b == '\f' || b == '\v' || b == 0x85 || b == 0xA0
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Float32To64(f32 float32) float64 {
	s := fmt.Sprintf("%v", f32)
	f64, _ := strconv.ParseFloat(s, 64)
	return f64
}
