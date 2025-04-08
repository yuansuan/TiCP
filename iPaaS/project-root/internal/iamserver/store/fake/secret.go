package fake

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type secrets struct {
	ds *datastore
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{ds: ds}
}

func (s *secrets) Create(ctx context.Context, secret *dao.Secret) error {
	s.ds.Lock()
	defer s.ds.Unlock()
	s.ds.secrets = append(s.ds.secrets, secret)
	return nil
}

func (s *secrets) Update(ctx context.Context, secret *dao.Secret) error {
	s.ds.Lock()
	defer s.ds.Unlock()
	for i, v := range s.ds.secrets {
		if v.AccessKeyId == secret.AccessKeyId {
			s.ds.secrets[i] = secret
			return nil
		}
	}
	return nil
}

func (s *secrets) Delete(ctx context.Context, akId string) error {
	s.ds.Lock()
	defer s.ds.Unlock()
	for i, v := range s.ds.secrets {
		if v.AccessKeyId == akId {
			s.ds.secrets = append(s.ds.secrets[:i], s.ds.secrets[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *secrets) DeleteByParentUser(ctx context.Context, akId, parentUser string) error {
	s.ds.Lock()
	defer s.ds.Unlock()
	for i, v := range s.ds.secrets {
		if v.AccessKeyId == akId && v.ParentUser == parentUser {
			s.ds.secrets = append(s.ds.secrets[:i], s.ds.secrets[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *secrets) List(ctx context.Context, parentUser string, offset, limit int) ([]*dao.Secret, error) {
	s.ds.RLock()
	defer s.ds.RUnlock()
	var secrets []*dao.Secret
	for _, v := range s.ds.secrets {
		if v.ParentUser == parentUser {
			secrets = append(secrets, v)
		}
	}
	return secrets, nil
}

func (s *secrets) ListAll(ctx context.Context, offset, limit int) ([]*dao.Secret, error) {
	s.ds.RLock()
	defer s.ds.RUnlock()
	return s.ds.secrets, nil
}

func (s *secrets) Get(ctx context.Context, akId string) (*dao.Secret, error) {
	s.ds.RLock()
	defer s.ds.RUnlock()
	for _, v := range s.ds.secrets {
		if v.AccessKeyId == akId {
			return v, nil
		}
	}
	return nil, nil
}

func (s *secrets) GetByUserID(ctx context.Context, userID string) (*dao.Secret, error) {
	s.ds.RLock()
	defer s.ds.RUnlock()
	for _, v := range s.ds.secrets {
		if v.ParentUser == userID {
			return v, nil
		}
	}
	return nil, nil
}

func (s *secrets) CleanExpireSecret(ctx context.Context, now time.Time) error {
	s.ds.Lock()
	defer s.ds.Unlock()
	for i, v := range s.ds.secrets {
		if v.Expiration.Before(now) {
			s.ds.secrets = append(s.ds.secrets[:i], s.ds.secrets[i+1:]...)
		}
	}
	return nil
}

func (s *secrets) UpdateTag(ctx context.Context, akID, userID, tag string) error {
	s.ds.Lock()
	defer s.ds.Unlock()
	for i, v := range s.ds.secrets {
		if v.AccessKeyId == akID && v.ParentUser == userID {
			s.ds.secrets[i].Tag = tag
			return nil
		}
	}
	return nil
}
