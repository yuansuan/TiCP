package sbs

import (
	"reflect"
	"unsafe"
)

// String 将字节数组转换为字符串类型
func String(bs []byte) (s string) {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = (*reflect.SliceHeader)(unsafe.Pointer(&bs)).Data
	hdr.Len = (*reflect.SliceHeader)(unsafe.Pointer(&bs)).Len
	return
}

// Bytes 将字符串转换为字节数组类型
func Bytes(s string) (bs []byte) {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	hdr.Len = (*reflect.StringHeader)(unsafe.Pointer(&s)).Len
	hdr.Cap = hdr.Len
	return
}
