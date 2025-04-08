package mysql

import (
	"context"
	"time"

	"github.com/marmotedu/errors"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"gorm.io/gorm"
)

type secrets struct {
	db *gorm.DB
}

func newSecrets(db *gorm.DB) *secrets {
	res := &secrets{db}
	return res
}

// Create creates a new secret.
func (s *secrets) Create(ctx context.Context, secret *dao.Secret) error {
	return s.db.Create(secret).Error
}

// Update updates an secret information by the secret identifier.
func (s *secrets) Update(ctx context.Context, secret *dao.Secret) error {
	return s.db.Save(secret).Error
}

// Delete deletes the secret by the secret identifier.
func (s *secrets) Delete(ctx context.Context, akId string) error {
	i := s.db.Where("accessKeyId= ?", akId).Delete(&dao.Secret{}).RowsAffected
	if i == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

// Delete deletes the secret by the secret identifier.
func (s *secrets) DeleteByParentUser(ctx context.Context, akID, parentUser string) error {
	i := s.db.Where("accessKeyId= ? and parentUser = ?", akID, parentUser).Delete(&dao.Secret{}).RowsAffected
	if i == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

// Get return an secret by the secret identifier.
func (s *secrets) Get(ctx context.Context, accessKey string) (*dao.Secret, error) {
	secret := dao.Secret{}
	err := s.db.Where("accessKeyId = ? ", accessKey).First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return &secret, nil
}

func (s *secrets) GetByUserID(ctx context.Context, userID string) (*dao.Secret, error) {
	secret := dao.Secret{}
	err := s.db.Where("parentUser = ? ", userID).First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return &secret, nil
}

// List return all secrets.
func (s *secrets) List(ctx context.Context, parentUser string, offset, limit int) ([]*dao.Secret, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	var ret []*dao.Secret
	d := s.db.Where("parentUser = ? ", parentUser).
		Offset(offset).
		Limit(limit).
		Find(&ret)
	return ret, d.Error
}

func (s *secrets) ListAll(ctx context.Context, offset, limit int) ([]*dao.Secret, error) {
	var ret []*dao.Secret
	d := s.db.Raw("select accessKeyId, accessKeySecret, sessionToken, parentUser, tag, expiration from secret limit ?,?", offset, limit).Scan(&ret)
	return ret, d.Error
}

func (s *secrets) UpdateTag(ctx context.Context, akID, userID, tag string) error {
	i := s.db.Where("accessKeyId = ? and parentUser = ?", akID, userID).Updates(&dao.Secret{Tag: tag}).RowsAffected
	if i == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

func (s *secrets) CleanExpireSecret(ctx context.Context, now time.Time) error {
	return s.db.Where("expiration < ?", now).Delete(&dao.Secret{}).Error
}
