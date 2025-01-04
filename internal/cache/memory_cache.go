package cache

import (
	"sync"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

type MemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

var _ Cache = (*MemoryCache)(nil)

func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration),
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found || time.Now().After(item.Expiration) {
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
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if time.Now().After(item.Expiration) {
				delete(c.items, key)
				logger.Debug("Cache item expired", 
					"key", key)
			}
		}
		c.mu.Unlock()
	}
}
