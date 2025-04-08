package v1

import (
	"context"

	"github.com/marmotedu/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/code"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type SecretSvc interface {
	CreateSecret(ctx context.Context, userID, tag string) (*dao.Secret, error)
	Get(ctx context.Context, akID string) (*dao.Secret, error)
	ListByParentUserID(ctx context.Context, userID string, offset, limit int) ([]*dao.Secret, error)
	List(ctx context.Context, offset, limit int) ([]*dao.Secret, error)
	DeleteByParentUser(ctx context.Context, akID, userID string) error

	AdminDelete(ctx context.Context, akID string) error

	UpdateTag(ctx context.Context, akID, userID, tag string) error

	GetByUserID(ctx context.Context, userID string) (*dao.Secret, error)
}

type secretService struct {
	store store.Factory
}

var _ SecretSvc = (*secretService)(nil)

func newSecrets(s *svc) *secretService {
	return &secretService{
		store: s.store,
	}
}

func (s *secretService) CreateSecret(ctx context.Context, userID, tag string) (*dao.Secret, error) {

	akID, adSecret, err := GenerateCredentials()
	if err != nil {
		return nil, err
	}
	secret := &dao.Secret{
		AccessKeyId:     akID,
		AccessKeySecret: adSecret,
		ParentUser:      userID,
		Tag:             tag,
		Expiration:      common.GetBigTime(),
	}
	err = s.store.Secrets().Create(ctx, secret)
	if err != nil {
		logging.Default().Infof("create secret error: %v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return secret, nil
}

func (s *secretService) Get(ctx context.Context, akID string) (*dao.Secret, error) {
	secret, err := s.store.Secrets().Get(ctx, akID)
	if err != nil {
		logging.Default().Infof("get secret error: %v", err)
		return nil, err
	}
	return secret, nil
}

func (s *secretService) GetByUserID(ctx context.Context, userID string) (*dao.Secret, error) {
	secret, err := s.store.Secrets().GetByUserID(ctx, userID)
	if err != nil {
		logging.Default().Infof("get secret by userID error: %v", err)
		return nil, err
	}
	return secret, nil
}

func (s *secretService) DeleteByParentUser(ctx context.Context, akID, userID string) error {
	err := s.store.Secrets().DeleteByParentUser(ctx, akID, userID)
	if err != nil {
		logging.Default().Infof("delete secret error: %v", err)
		return err
	}
	return nil
}

func (s *secretService) List(ctx context.Context, offset, limit int) ([]*dao.Secret, error) {
	secrets, err := s.store.Secrets().ListAll(ctx, offset, limit)
	if err != nil {
		logging.Default().Infof("list secret error: %v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return secrets, nil
}

func (s *secretService) ListByParentUserID(ctx context.Context, userID string, offset, limit int) ([]*dao.Secret, error) {
	secrets, err := s.store.Secrets().List(ctx, userID, offset, limit)
	if err != nil {
		logging.Default().Infof("list secret error: %v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return secrets, nil
}

func (s *secretService) AdminDelete(ctx context.Context, akID string) error {
	err := s.store.Secrets().Delete(ctx, akID)
	if err != nil {
		logging.Default().Infof("delete secret error: %v", err)
		return err
	}
	return nil
}

func (s *secretService) UpdateTag(ctx context.Context, akID, userID, tag string) error {
	err := s.store.Secrets().UpdateTag(ctx, akID, userID, tag)
	if err != nil {
		logging.Default().Infof("update secret tag error: %v", err)
		return err
	}
	return nil
}
