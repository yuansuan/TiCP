package utils

import "encoding/base64"

// StringPtr 获取一个字符串的指针
func StringPtr(s string) *string {
	return &s
}

// Base64StringPtr 使用Base64编码字符串后返回指针
func Base64StringPtr(s string) *string {
	return StringPtr(base64.StdEncoding.EncodeToString([]byte(s)))
}

func BoolPtr(b bool) *bool {
	return &b
}
