package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

// DES key length must be 8
const secretKey = "!@#$%^&*"

// Encrypt Encrypt
func Encrypt(plainText string) (string, error) {
	src := []byte(plainText)
	key := []byte(secretKey)

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	dst := make([]byte, len(src))
	// encrypt plaintext
	blockMode.CryptBlocks(dst, src)
	// convert byte array into string
	cipherText := base64.StdEncoding.EncodeToString(dst)
	return cipherText, nil
}

// Decrypt Decrypt
func Decrypt(cipherText string) (string, error) {
	key := []byte(secretKey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, key)
	src, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	dst := make([]byte, len(src))
	// decrypt cipherText
	blockMode.CryptBlocks(dst, src)
	dst = unPadding(dst)
	return string(dst), nil
}

// PKCS5Padding
func padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padText...)
}

// Remove PKCS5Padding
func unPadding(cipherText []byte) []byte {
	length := len(cipherText)
	unPadNum := int(cipherText[length-1])
	return cipherText[:(length - unPadNum)]
}
