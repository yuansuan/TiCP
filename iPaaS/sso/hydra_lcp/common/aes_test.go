package common

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAES(t *testing.T) {
	key := "111111111111111111111111"
	src := "0000"
	assert := require.New(t)

	cSrc, err := AESEncrypt([]byte(src), []byte(key))

	base64CSrc := base64.StdEncoding.EncodeToString(cSrc)

	cipherPwd, _ := base64.StdEncoding.DecodeString(base64CSrc)

	dSrc, err := AESDecrypt(cipherPwd, []byte(key))

	fmt.Println(base64CSrc, cSrc, err, len(cipherPwd), cipherPwd)

	assert.Equal(src, string(dSrc))
}
