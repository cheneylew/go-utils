package utils

import (
	"unsafe"
	"fmt"
)

// BytesToString convert []byte type to string type.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes convert string type to []byte type.
// NOTE: panic if modify the member value of the []byte.
func StringToBytes(s string) []byte {
	sp := *(*[2]uintptr)(unsafe.Pointer(&s))
	bp := [3]uintptr{sp[0], sp[1], sp[1]}
	return *(*[]byte)(unsafe.Pointer(&bp))
}

func FormatFlow(byteCount float64) string {
	if byteCount < 1024 {
		return fmt.Sprintf("%.2f byte", byteCount)
	} else if byteCount < 1024*1024 {
		return fmt.Sprintf("%.2f KB", byteCount/1024)
	} else if byteCount < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", byteCount/(1024 * 1024))
	} else if byteCount < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", byteCount/(1024 * 1024*1024))
	}

	return ""
}