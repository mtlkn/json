package json

import "unsafe"

func stringToBytes(s string) []byte {
	if s == "" {
		return []byte{}
	}

	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return *(*string)(unsafe.Pointer(&b))
}
