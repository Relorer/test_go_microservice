package cache

import (
	"sync"
	"time"
)

type SimpleCache struct {
	data       map[string]interface{}
	expiration map[string]time.Time
	mu         sync.RWMutex
	ttl        time.Duration
}

func NewSimpleCache(ttl, cleanupInterval time.Duration) *SimpleCache {
	c := &SimpleCache{
		data:       make(map[string]interface{}),
		expiration: make(map[string]time.Time),
		ttl:        ttl,
	}
	go c.cleanup(cleanupInterval)
	return c
}

func (c *SimpleCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiration[key] = time.Now().Add(c.ttl)
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if !ok {
		return nil, false
	}

	expiration, ok := c.expiration[key]
	if !ok {
		return value, true
	}

	if expiration.After(time.Now()) {
		return value, true
	}

	// Delete outdated data
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.data, key)
		delete(c.expiration, key)
	}()

	return nil, false
}

func (c *SimpleCache) cleanup(interval time.Duration) {
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
