package cache

import (
	"time"
)

// Cache defines a generic caching interface
type Cache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}
	go cache.cleanupLoop()
	return cache
}
