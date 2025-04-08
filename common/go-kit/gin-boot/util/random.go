// Copyright (C) 2019 LambdaCal Inc.

package util

import (
	"math/rand"
	"time"
)

const constCharCollections = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomData gen random data len by byte num
func RandomData(byteNum int) ([]byte, error) {
	data := make([]byte, byteNum)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}

// RandomString is get random string
func RandomString(strLen int) string {
	bytes := []byte(constCharCollections)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < strLen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
