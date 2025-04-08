//go:build darwin
// +build darwin

package dao

import (
	"log"
	"os"
	"testing"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"golang.org/x/net/context"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
)

var d = UserDao{}

func TestMain(m *testing.M) {
	os.Chdir("..")
	_ = boot.Default()
	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	ctx := context.TODO()
	var user *models.SsoUser

	t.Run("add a user", func(t *testing.T) {
		user = &models.SsoUser{
			Ysid:  12,
			Phone: "12",
			Email: "test1@yuansuan.cn",
		}
	})

	err := d.Add(ctx, user)
	log.Print(err)
}

func TestUpdate(t *testing.T) {
	ctx := context.TODO()
	var user models.SsoUser

	t.Run("add email", func(t *testing.T) {
		user = models.SsoUser{
			Ysid:  124,
			Email: "test@yuansuan.cn",
		}
	})

	err := d.Update(ctx, user)
	log.Print(err)
}

func TestGet(t *testing.T) {
	ctx := context.TODO()
	var user models.SsoUser

	t.Run("get user", func(t *testing.T) {
		user.Email = "test@yuansuan.cn"
		user.PwdHash = ""
	})

	_, err := d.Get(ctx, &user)
	log.Println(user)
	log.Println(err)

}
