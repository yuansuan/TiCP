package strutil

import (
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 随机大写字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := random.Intn(26) + 65
		bytes[i] = byte(b)
	}

	return string(bytes)
}
