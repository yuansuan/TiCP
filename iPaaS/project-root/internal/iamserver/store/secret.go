package store

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type SecretStore interface {
	Create(ctx context.Context, secret *dao.Secret) error
	Update(ctx context.Context, secret *dao.Secret) error
	Delete(ctx context.Context, akId string) error

	DeleteByParentUser(ctx context.Context, akID, parentUser string) error
	List(ctx context.Context, parentUser string, offset, limit int) ([]*dao.Secret, error)

	ListAll(ctx context.Context, offset, limit int) ([]*dao.Secret, error)
	Get(ctx context.Context, akId string) (*dao.Secret, error)
	CleanExpireSecret(ctx context.Context, now time.Time) error

	GetByUserID(ctx context.Context, userID string) (*dao.Secret, error)

	UpdateTag(ctx context.Context, akID, userID, tag string) error
}
