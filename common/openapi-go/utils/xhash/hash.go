package xhash

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"

	"github.com/yuansuan/ticp/common/openapi-go/utils"
)

// MD5 returns the md5 hash of the text, in hexadecimal
func MD5(text string) string {
	h := md5.New()
	h.Write(utils.Bytes(text))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 returns the sha1 hash of the text, in hexadecimal
func SHA1(text string) string {
	h := sha1.New()
	h.Write(utils.Bytes(text))
	return hex.EncodeToString(h.Sum(nil))
}
