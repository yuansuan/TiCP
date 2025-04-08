package cache

import (
	"errors"
	"sync"
	"time"

	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
)

type Cache struct {
	m    sync.Mutex
	data map[string]*cacheItem
}

type cacheItem struct {
	secret *iam_api.CacheSecret
	expire time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*cacheItem),
	}
}

func (c *Cache) GetSecret(accessKeyId string) (*iam_api.CacheSecret, error) {
	c.m.Lock()
	defer c.m.Unlock()

	item, ok := c.data[accessKeyId]
	if !ok {
		return nil, ErrSecretNotFound
	}
	if item.expire.Before(time.Now()) {
		delete(c.data, accessKeyId)
		return nil, ErrSecretNotFound
	}
	return item.secret, nil
}

func (c *Cache) SetSecret(key string, secret *iam_api.CacheSecret) {
	item := &cacheItem{
		secret: secret,
		expire: secret.Expire,
	}
	c.data[key] = item
}

func (c *Cache) ClearAndSet(secrets []*iam_api.CacheSecret) {
	c.m.Lock()
	defer c.m.Unlock()
	c.data = make(map[string]*cacheItem)
	for _, secret := range secrets {
		c.SetSecret(secret.AccessKeyId, secret)
	}
}
