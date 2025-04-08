package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache Cache
var Cache *cache.Cache

func init() {
	Cache = cache.New(5*time.Minute, 15*time.Minute)
}
