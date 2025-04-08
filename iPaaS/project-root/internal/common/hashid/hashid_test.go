package hashid

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

const (
	TestHashKey = "00010203040506070809aabbccddeeff"
	TestHashIv  = "1122334455667788990a0b0c0d0e0f00"

	TestSnowflakeId snowflake.ID = 123456789987654321
)

func TestNew(t *testing.T) {
	c, err := New(TestHashKey, TestHashIv)

	assert.NotNil(t, c)
	assert.NoError(t, err)
}

func TestCodec_Encode(t *testing.T) {
	c, err := New(TestHashKey, TestHashIv)
	if assert.NoError(t, err) && assert.NotNil(t, c) {
		for i := 0; i < 32; i++ {
			id, err := c.Encode(TestSnowflakeId)
			if assert.NoError(t, err) && assert.NotEmpty(t, id) {
				t.Logf("Encode: %d -> %s", TestSnowflakeId, id)
			}
		}
	}
}

func TestCodec_Decode(t *testing.T) {
	c, err := New(TestHashKey, TestHashIv)
	if assert.NoError(t, err) && assert.NotNil(t, c) {
		for i := 0; i < 32; i++ {
			id, err := c.Encode(TestSnowflakeId)
			if assert.NoError(t, err) && assert.NotEmpty(t, id) {
				t.Logf("Encode: %d -> %s", TestSnowflakeId, id)
				if decodeId, err := c.Decode(id); assert.NoError(t, err) {
					assert.Equal(t, TestSnowflakeId, decodeId)
					t.Logf("Decode: %s -> %d", id, decodeId)
				}
			}
		}
	}
}

func TestCodecStr_DecodeStr(t *testing.T) {
	c, err := New(TestHashKey, TestHashIv)
	if assert.NoError(t, err) && assert.NotNil(t, c) {
		for i := 0; i < 100; i++ {
			inputStr := genRandStr(rand.Intn(9999) + 1)
			encodedStr, err := c.EncodeStr(inputStr)
			if assert.NoError(t, err) && assert.NotEmpty(t, encodedStr) {
				//t.Logf("Encode: %s -> %s", inputStr, encodedStr)
				if decodedStr, err := c.DecodeStr(encodedStr); assert.NoError(t, err) {
					assert.Equal(t, inputStr, decodedStr)
					//t.Logf("Decode: %s -> %s", inputStr, decodedStr)
				}
			}
		}
	}
}

func genRandStr(n int) string {
	res := strings.Builder{}
	for i := 0; i < n; i++ {
		cInt := rand.Intn(256)
		res.WriteByte(byte(cInt))
	}
	return res.String()
}
