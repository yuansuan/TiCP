package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//参考文档
//http://www.topgoer.com/%E5%85%B6%E4%BB%96/%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86/%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86.html
//高级加密标准（Adevanced Encryption Standard ,AES）

//16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法

// PKCS7Padding  PKCS7补位字符
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0x00}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 解密时删除补位字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)

	// 长度正好是aes.BlockSize的整倍数
	if len(origData)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("cipher text is not a multiple of the block size, length: %v", length)
	}

	return bytes.TrimRight(origData, string([]byte{0x00})), nil
}

var iv = []byte{
	0x30, 0x30, 0x30, 0x30,
	0x30, 0x30, 0x30, 0x30,
	0x30, 0x30, 0x30, 0x30,
	0x30, 0x30, 0x30, 0x30,
}

// AESEncrypt 加密
func AESEncrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AESDecrypt 解密
func AESDecrypt(cypted []byte, key []byte) (string, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}
