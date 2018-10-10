package cache

import "time"

// Cache represents an
type Cache struct {
	values map[string]interface{}

	capacity int64
	expiry   time.Duration
}

// New creates a new LRU cache with a global expiry
func New(capacity int64, expiry time.Duration) *Cache {
	return &Cache{
		values:   make(map[string]interface{}),
		capacity: capacity,
		expiry:   expiry,
	}
}

// Get returns the value for an associated key if it exists,
// otherwise it will nil.
func (c *Cache) Get(key string) (interface{}, bool) {
	if val, found := c.values[key]; found {
		return val, true
	}

	return nil, false
}

// Set creates a new record in the cache.s
func (c *Cache) Set(key string, val interface{}) (interface{}, bool) {
	c.values[key] = val

	return nil, false
}
