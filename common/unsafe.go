package common

import (
	"reflect"
	"unsafe"
)

// BytesToString converts bytes to string without recursion
func BytesToString(b []byte) string {
	var (
		bh = (*reflect.SliceHeader)(unsafe.Pointer(&b))
		sh = reflect.StringHeader{bh.Data, bh.Len}
	)
	return *(*string)(unsafe.Pointer(&sh))
}

// StringToBytes converts string to bytes without recursion
func StringToBytes(s string) []byte {
	var (
		sh = (*reflect.StringHeader)(unsafe.Pointer(&s))
		bh = reflect.SliceHeader{sh.Data, sh.Len, 0}
	)
	return *(*[]byte)(unsafe.Pointer(&bh))
}
