package strutil

import (
	"bytes"
	"encoding/base64"
	"io"
)

// Base64Decode base64字符串解码
func Base64Decode(s string) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(s))
	return io.ReadAll(decoder)
}
