package cache

import (
	"sync"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

const (
	DefaultExpiration = 10 * time.Minute
	CleanupInterval   = 5 * time.Minute
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type MemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if duration <= 0 {
		duration = DefaultExpiration
	}

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(CleanupInterval)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *MemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range c.items {
		if now > v.Expiration {
			delete(c.items, k)
		}
	}

	logger.Info("Cache cleanup completed", 
		"expired_items_removed", len(c.items))
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}
	go cache.cleanupLoop()
	return cache
}
