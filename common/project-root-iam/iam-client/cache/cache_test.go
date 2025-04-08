package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

type Suite struct {
	suite.Suite
	secrets []*iam_api.CacheSecret
	cache   *Cache
}

func (s *Suite) SetupSuite() {
	secret1 := &iam_api.CacheSecret{
		AccessKeyId:     "access_key_id_1",
		AccessKeySecret: "access_key_secret_1",
		// after 5 minutes
		Expire: time.Now().Add(5 * time.Minute),
	}
	secret2 := &iam_api.CacheSecret{
		AccessKeyId:     "access_key_id_2",
		AccessKeySecret: "access_key_secret_2",
		// before 5 minutes
		Expire: time.Now().Add(-5 * time.Minute),
	}
	s.secrets = []*iam_api.CacheSecret{secret1, secret2}

	s.cache = NewCache()
	for _, secret := range s.secrets {
		s.cache.SetSecret(secret.AccessKeyId, secret)
	}
}

func TestCache(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_GetSecret() {
	// len of cache equals 2
	s.NotNil(s.cache.data[s.secrets[0].AccessKeyId])
	s.NotNil(s.cache.data[s.secrets[1].AccessKeyId])

	// except get secret1
	secret, err := s.cache.GetSecret(s.secrets[0].AccessKeyId)
	s.Nil(err)
	s.Equal(s.secrets[0], secret)

	// get secret2 failed, error equals ErrSecretNotFound
	_, err = s.cache.GetSecret(s.secrets[1].AccessKeyId)
	s.Equal(ErrSecretNotFound, err)

	// len of cache equals 1
	s.NotNil(s.cache.data[s.secrets[0].AccessKeyId])
	s.Nil(s.cache.data[s.secrets[1].AccessKeyId])
}

func (s *Suite) Test_ClearAndSetSecret() {
	// len of cache equals 2
	s.NotNil(s.cache.data[s.secrets[0].AccessKeyId])
	s.NotNil(s.cache.data[s.secrets[1].AccessKeyId])

	secret3 := &iam_api.CacheSecret{
		AccessKeyId:     "access_key_id_3",
		AccessKeySecret: "access_key_secret_3",
		// after 5 minutes
		Expire: time.Now().Add(5 * time.Minute),
	}

	s.secrets = append(s.secrets, secret3)
	s.cache.ClearAndSet(s.secrets)

	// len of cache equals 3
	s.NotNil(s.cache.data[s.secrets[0].AccessKeyId])
	s.NotNil(s.cache.data[s.secrets[1].AccessKeyId])
	s.NotNil(s.cache.data[s.secrets[2].AccessKeyId])
}
