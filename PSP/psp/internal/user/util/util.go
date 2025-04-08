package util

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// PasswdCrypto crypto passwd into hash
// use sha256
func PasswdCrypto(passwd string) string {
	crypto := sha256.New()
	crypto.Write([]byte(passwd))
	return fmt.Sprintf("%x", crypto.Sum(nil))
}

func GetClientIPAndHost(ctx context.Context) string {

	logger := logging.GetLogger(ctx)

	ip, err := util.GetInMetadata(ctx, "x-real-ip")
	if err != nil {
		logger.Errorf("Failed to get the X-Real-IP value, %v", err.Error())
		return ""
	}

	logger.Infof("The X-Real-IP value is %v.", ip)

	return ip
}

// 生成随机密码
func GenerateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}
