package conf_type

// TypeCacheRedis TypeCacheRedis
const TypeCacheRedis = "redis"

// Cache cache
type Cache struct {
	BackendType string `yaml:"backend_type"`
	Name        string `yaml:"name"`
}

// Caches Caches
type Caches map[string]Cache
