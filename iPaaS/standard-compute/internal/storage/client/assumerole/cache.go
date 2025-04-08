package assumerole

import (
	"sync"
	"time"
)

type Value struct {
	AccessKeyId     string
	AccessKeySecret string
	Token           string
	ExpiredTime     time.Time
}

// Cache 当前用户不会很多，此处map也不会膨胀
type Cache struct {
	m *sync.Map
}

func NewCache() *Cache {
	return &Cache{
		m: &sync.Map{},
	}
}

func (c *Cache) Get(userId string) (Value, bool) {
	v, exist := c.m.Load(userId)
	if !exist {
		return Value{}, false
	}

	return v.(Value), true
}

func (c *Cache) Set(userId string, v Value) {
	c.m.Store(userId, v)
}
