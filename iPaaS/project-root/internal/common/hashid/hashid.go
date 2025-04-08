package hashid

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"time"
)

const (
	InvalidHashId snowflake.ID = 0

	defaultHashKey = "deb516af9a691b81c8662a02a57ac974"
	defaultHashIv  = "55322c3159736842ce296c0864f916e9"
)

var (
	defaultCodec         *Codec
	ErrHashBlockTooShort = errors.New("block too short")
	ErrHashBlockSize     = errors.New("block size not available")
)

func init() {
	defaultCodec, _ = New(defaultHashKey, defaultHashIv)
}

type Codec struct {
	key, iv []byte

	block cipher.Block
}

func (c *Codec) Encode(id snowflake.ID) (r string, err error) {
	defer func() {
		if rErr := recover(); rErr != nil {
			id = InvalidHashId
			err = errors.New(fmt.Sprintf("%s", rErr))
		}
	}()

	src := make([]byte, 0, aes.BlockSize)
	src = append(src, decodeInt64(int64(id))...)
	src = append(src, decodeInt64(time.Now().UnixNano())...)
	dst := make([]byte, aes.BlockSize)

	cbc := cipher.NewCBCEncrypter(c.block, c.iv)
	cbc.CryptBlocks(dst, src)

	return hex.EncodeToString(dst), nil
}

func (c *Codec) Decode(s string) (id snowflake.ID, err error) {
	defer func() {
		if rErr := recover(); rErr != nil {
			id = InvalidHashId
			err = errors.New(fmt.Sprintf("%s", rErr))
		}
	}()

	data, err := hex.DecodeString(s)
	if err != nil {
		return InvalidHashId, errors.Wrap(err, "hex")
	}
	if len(data) != aes.BlockSize {
		return InvalidHashId, ErrHashBlockTooShort
	}

	dst := make([]byte, aes.BlockSize)
	cbc := cipher.NewCBCDecrypter(c.block, c.iv)
	cbc.CryptBlocks(dst, data)

	return snowflake.ID(encodeInt64(dst[:8])), nil
}

func (c *Codec) EncodeStr(s string) (r string, err error) {
	defer func() {
		if rErr := recover(); rErr != nil {
			err = errors.New(fmt.Sprintf("%s", rErr))
		}
	}()
	padding := PKCS7Padding([]byte(s), aes.BlockSize)
	dst := make([]byte, len(padding))

	cbc := cipher.NewCBCEncrypter(c.block, c.iv)
	cbc.CryptBlocks(dst, padding)

	return hex.EncodeToString(dst), nil
}

func (c *Codec) DecodeStr(s string) (res string, err error) {
	defer func() {
		if rErr := recover(); rErr != nil {
			err = errors.New(fmt.Sprintf("%s", rErr))
		}
	}()

	data, err := hex.DecodeString(s)
	if err != nil {
		return "", errors.Wrap(err, "hex")
	}
	if len(data)%aes.BlockSize != 0 {
		return "", ErrHashBlockSize
	}

	dst := make([]byte, len(data))
	cbc := cipher.NewCBCDecrypter(c.block, c.iv)
	cbc.CryptBlocks(dst, data)
	resBytes, err := PKCS7UnPadding(dst)
	return string(resBytes), err
}

func New(key, iv string) (*Codec, error) {
	k, err := decodeBlockHex(key)
	if err != nil {
		return nil, errors.Wrap(err, "key")
	}
	i, err := decodeBlockHex(iv)
	if err != nil {
		return nil, errors.Wrap(err, "iv")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, errors.Wrap(err, "aes")
	}

	return &Codec{key: k, iv: i, block: block}, nil
}

func Encode(id snowflake.ID) (r string, err error) {
	return defaultCodec.Encode(id)
}

func Decode(s string) (id snowflake.ID, err error) {
	return defaultCodec.Decode(s)
}

func EncodeStr(s string) (string, error) {
	return defaultCodec.EncodeStr(s)
}

func DecodeStr(s string) (string, error) {
	return defaultCodec.DecodeStr(s)
}

func decodeBlockHex(s string) ([]byte, error) {
	bs, err := hex.DecodeString(s)
	if err != nil {
		return nil, errors.Wrap(err, "hex")
	}
	if len(bs) != aes.BlockSize {
		return nil, ErrHashBlockTooShort
	}

	return bs, nil
}

func decodeInt64(v int64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(v))
	return data
}

func encodeInt64(v []byte) int64 {
	return int64(binary.BigEndian.Uint64(v))
}

// PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反向操作，删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("empty input data")
	} else {
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)], nil
	}
}
