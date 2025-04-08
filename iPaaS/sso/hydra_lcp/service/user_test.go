//go:build darwin
// +build darwin

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
)

var userSrv = NewUserSrv()

func TestMain(m *testing.M) {
	os.Chdir("..")
	_ = boot.Default()
	os.Exit(m.Run())
}

func TestUserService_AddUser(t *testing.T) {
	ctx := context.TODO()
	userID := int64(2)
	email := "fwchen@yuansuan.cn"
	pwd := "helloworld"
	u := &models.SsoUser{Ysid: userID, Email: email}

	err := userSrv.Add(ctx, u, pwd)
	fmt.Println(err)
}
