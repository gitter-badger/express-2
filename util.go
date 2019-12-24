package express

import (
	"reflect"
	"unsafe"
)

// Credits to @savsgio
// https://github.com/savsgio/gotils/blob/master/conv.go

// b2s converts byte slice to a string without memory allocation.
func b2s(b []byte) string {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&b))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*string)(unsafe.Pointer(&bh))
}

// s2b converts string to a byte slice without memory allocation.
func s2b(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
