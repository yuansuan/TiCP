// Copyright (C) 2019 LambdaCal Inc.

package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// AESEncrypt byte to secret byte
// return data length may be different with input, be careful
func AESEncrypt(origData, iv, key []byte) ([]byte, error) {
	if len(origData) == 0 {
		return nil, fmt.Errorf("origData is empty")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if block.BlockSize() != len(iv) {
		return nil, fmt.Errorf("iv size %v not equal block size %v", len(iv), block.BlockSize())
	}

	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypt := make([]byte, len(origData))
	blockMode.CryptBlocks(crypt, origData)
	return crypt, nil
}

// AESDecrypt secret byte to byte
func AESDecrypt(crypt, iv, key []byte) ([]byte, error) {
	if len(crypt) == 0 {
		return nil, fmt.Errorf("crypt is empty")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if block.BlockSize() != len(iv) {
		return nil, fmt.Errorf("iv size %v not equal block size %v", len(iv), block.BlockSize())
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypt))
	blockMode.CryptBlocks(origData, crypt)
	origData = PKCS5UnPadding(origData)
	if origData == nil {
		return nil, fmt.Errorf("secret decode error")
	}
	return origData, nil
}

// PKCS5Padding is padding byte
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

// PKCS5UnPadding is un padding byte
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	if unPadding >= length {
		return nil
	}
	return origData[:(length - unPadding)]
}

// Base64Encode is Base64Encode
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// Base64Decode is Base64Decode
func Base64Decode(input string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HMACSHA1Encrypt encrype data as hmac-sha1
func HMACSHA1Encrypt(origData, key []byte) []byte {
	return hmac.New(sha1.New, key).Sum(origData)
}

// MD5 sum input string's md5
func MD5(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
