package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCache(t *testing.T) {
	cache := NewSimpleCache(time.Second, time.Second*2)

	// Test Set and Get methods
	cache.Set("key1", "value1")
	cache.Set("key2", 42)
	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)
	value, ok = cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, 42, value)

	// Test Get method with non-existent key
	_, ok = cache.Get("key3")
	assert.False(t, ok)

	// Test Get method with expired key
	cache.Set("key4", "value4")
	time.Sleep(time.Second)
	_, ok = cache.Get("key4")
	assert.False(t, ok)
}
