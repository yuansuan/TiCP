/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package middleware

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// DefaultExpireDur DefaultExpireDur
var DefaultExpireDur = time.Minute * 5

// ICache ICache
type ICache interface {
	Put(prefix, key string, value interface{}) error
	PutWithExpire(prefix, key string, value interface{}, expire time.Duration) error
	PutUnExpired(prefix, key string, value interface{}) error
	Get(prefix, key string, vtype interface{}) (interface{}, bool)
	Delete(prefix, key string)
}

// RedisCache RedisCache
type RedisCache struct {
	backend *redis.Client
}

// NewRedisCache NewRedisCache
func NewRedisCache(backend *redis.Client) ICache {
	return &RedisCache{backend: backend}
}

func generate(prefix, key string) string {
	return prefix + "-" + key
}

func (c *RedisCache) put(prefix, key string, value interface{}, expire time.Duration) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := c.backend.Set(generate(prefix, key), string(content), expire)
	_, err = cmd.Result()
	return err
}

// Put Put
func (c *RedisCache) Put(prefix, key string, value interface{}) error {
	return c.put(prefix, key, value, DefaultExpireDur)
}

// PutWithExpire PutWithExpire
func (c *RedisCache) PutWithExpire(prefix, key string, value interface{}, expire time.Duration) error {
	return c.put(prefix, key, value, expire)
}

// PutUnExpired PutUnExpired
func (c *RedisCache) PutUnExpired(prefix, key string, value interface{}) error {
	return c.put(prefix, key, value, 0)
}

// Get Get
func (c *RedisCache) Get(prefix, key string, vtype interface{}) (interface{}, bool) {
	cmd := c.backend.Get(generate(prefix, key))
	content := cmd.Val()

	if cmd.Err() != nil {
		return vtype, false
	}

	if err := json.Unmarshal([]byte(content), vtype); err != nil {
		logging.Default().Warnf("parse vtype failed %v", err)
		return vtype, false
	}
	return vtype, true
}

// Delete Delete
func (c *RedisCache) Delete(prefix, key string) {
	c.backend.Del(generate(prefix, key))
}
