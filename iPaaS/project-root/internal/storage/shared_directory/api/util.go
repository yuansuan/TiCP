package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
)

// ToResponseSharedDirectorys 转换为响应
func ToResponseSharedDirectorys(sharedDirectorys []*model.SharedDirectory) []*schema.SharedDirectory {
	res := make([]*schema.SharedDirectory, 0, len(sharedDirectorys))
	for _, sharedDirectory := range sharedDirectorys {
		res = append(res, ToResponseSharedDirectory(sharedDirectory))
	}
	return res
}

// ToResponseSharedDirectory 转换为响应
func ToResponseSharedDirectory(sharedDirectory *model.SharedDirectory) *schema.SharedDirectory {
	return &schema.SharedDirectory{
		Path:       sharedDirectory.Path,
		UserName:   addPrefixToUserName(sharedDirectory.SharedUserName),
		Password:   sharedDirectory.SharedPassword,
		SharedHost: sharedDirectory.SharedHost,
		SharedSrc:  sharedDirectory.SharedSrc,
	}
}

// DefaultPWLength 默认密码长度
const DefaultPWLength = 16

// UserPrefix 用户名前缀
const UserPrefix = "YS"

func generateUserName() (string, error) {
	idGen, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}

	id := idGen.Generate()

	username, err := hashid.Encode(id)
	if err != nil {
		return "", err
	}

	return username, nil
}

func addPrefixToUserName(username string) string {
	return fmt.Sprintf("%s%s", UserPrefix, username)
}

func generateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

func generateRandomPassword(length int) (string, error) {
	randomBytes, err := generateRandomBytes(length)
	if err != nil {
		return "", err
	}
	password := base64.URLEncoding.EncodeToString(randomBytes)

	// 截取指定长度的密码
	return password[:length], nil
}
