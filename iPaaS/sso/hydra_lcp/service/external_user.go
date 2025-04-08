package service

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
)

type externalUserService struct{}

// ExternalUserService ...
var ExternalUserService = &externalUserService{}

// Add Add
func (s *externalUserService) Add(ctx context.Context, user *models.SsoExternalUser) (err error) {
	err = dao.ExternalUserDao.Add(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// Get Get
func (s *externalUserService) Get(ctx context.Context, u *models.SsoExternalUser) error {
	return dao.ExternalUserDao.Get(ctx, u)
}

// GetID GetID
func (s *externalUserService) GetID(ctx context.Context, userName string) (userID int64, err error) {
	u := models.SsoExternalUser{UserName: userName}

	err = dao.ExternalUserDao.Get(ctx, &u)

	return u.Ysid, err
}
