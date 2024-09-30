package str

import (
	"reflect"
	"unsafe"
)

func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToByte(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
