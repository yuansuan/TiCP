package dao

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
)

type externalUserDao struct{}

// ExternalUserDao ExternalUserDao
var ExternalUserDao = &externalUserDao{}

// Add Add
func (d *externalUserDao) Add(ctx context.Context, user *models.SsoExternalUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}
	user.CreateTime = time.Now()
	_, err := session.Insert(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return nil
}

// Get Get
func (d *externalUserDao) Get(ctx context.Context, user *models.SsoExternalUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Get(user)
	if err != nil {
		return status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get user, err: %v", err.Error())
	}

	return nil
}
