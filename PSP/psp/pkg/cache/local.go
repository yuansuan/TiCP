package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// LocalCache ...
type LocalCache interface {
	Put(prefix, key string, value interface{})
	PutUnExpired(prefix, key string, value interface{})
	Get(prefix, key string) (interface{}, bool)
	Delete(prefix, key string)
}

type impl struct {
	cache *cache.Cache
}

// NewLocalCache ...
func NewLocalCache() LocalCache {
	return &impl{cache: cache.New(5*time.Minute, 10*time.Minute)}
}

// PutUnExpired ...
func (c *impl) PutUnExpired(prefix, key string, value interface{}) {
	c.cache.Set(prefix+"-"+key, value, cache.NoExpiration)
}

// Put ...
func (c *impl) Put(prefix, key string, value interface{}) {
	c.cache.Set(prefix+"-"+key, value, cache.DefaultExpiration)
}

// Get ...
func (c *impl) Get(prefix, key string) (interface{}, bool) {
	return c.cache.Get(prefix + "-" + key)
}

// Delete Delete
func (c *impl) Delete(prefix, key string) {
	c.cache.Delete(prefix + "-" + key)
}
