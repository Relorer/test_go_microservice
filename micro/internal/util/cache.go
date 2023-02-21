package util

import (
	"sync"
	"time"
)

type Cache struct {
	data       map[string]interface{}
	expiration map[string]time.Time
	mu         sync.RWMutex
	ttl        time.Duration
}

func NewCache(ttl, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		data:       make(map[string]interface{}),
		expiration: make(map[string]time.Time),
		ttl:        ttl,
	}
	go c.cleanup(cleanupInterval)
	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiration[key] = time.Now().Add(c.ttl)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key]; ok {
		if expiration, ok := c.expiration[key]; ok {
			if expiration.After(time.Now()) {
				return value, true
			} else {
				delete(c.data, key)
				delete(c.expiration, key)
			}
		} else {
			return value, true
		}
	}
	return nil, false
}

func (c *Cache) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.mu.Lock()
		for key, expiration := range c.expiration {
			if expiration.Before(time.Now()) {
				delete(c.data, key)
				delete(c.expiration, key)
			}
		}
		c.mu.Unlock()
	}
}
