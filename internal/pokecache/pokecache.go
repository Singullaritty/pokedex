package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entry map[string]cacheEntry
	mu    sync.RWMutex
}

type cacheEntry struct {
	CreatedAt time.Time
	Value     []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Entry: make(map[string]cacheEntry),
	}

	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			for _, v := range cache.Entry {
				elapsed := time.Since(v.CreatedAt)
				if interval < elapsed {
					cache.reapLoop()
				}
			}
		}

	}()

	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{CreatedAt: time.Now(), Value: value}
	c.Entry[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.Entry[key]
	return data.Value, ok
}

func (c *Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.Entry {
		delete(c.Entry, k)
	}
}
