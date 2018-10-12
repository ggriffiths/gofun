package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

// Cache represents an LRU cache with expiry cleanup
type Cache struct {
	lruCache        *lru.Cache
	expiry          time.Duration
	cleanupInterval time.Duration

	done chan bool
}

// Entry has the value and insertion timestamp
type Entry struct {
	InsertionTime time.Time
	Val           interface{}
}

// New creates a new LRU cache with a global expiry
func New(capacity int, expiry time.Duration, cleanupInterval time.Duration) (*Cache, error) {
	lc, err := lru.New(capacity)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create LRU Cache")
	}

	c := &Cache{
		lruCache:        lc,
		expiry:          expiry,
		cleanupInterval: cleanupInterval,
	}
	go c.cleanupLoop()

	return c, nil
}

// Close stops the cleanup goroutine
func (c *Cache) Close() {
	c.lruCache.Purge()
	c.done <- true
}

// Get returns the value for an associated key if it exists,
// otherwise it will nil.
func (c *Cache) Get(key string) (interface{}, bool) {
	if entry, found := c.lruCache.Get(key); found {
		log.Debug("Cache hit", "key", key)
		entryTyped := entry.(Entry)
		return entryTyped.Val, true
	}

	log.Debug("Cache miss", "key", key)
	return nil, false
}

// Set creates a new record in the cache.
func (c *Cache) Set(key string, val interface{}) (interface{}, bool) {
	c.lruCache.Add(key, Entry{time.Now(), val})

	return nil, false
}

// cleanupLoop runs in a separate goroutine to clean
// expired keys for a configurable interval.
func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)

	for {
		select {
		case <-ticker.C:
			keys := c.lruCache.Keys()

			// Remove expired keys starting from the oldest
			c.cleanup(keys)

		case <-c.done:
			return
		}
	}
}

// cleanup removes all expired keys
func (c *Cache) cleanup(keys []interface{}) {
	for _, k := range keys {
		if entry, found := c.lruCache.Get(k); found {
			entryTyped := entry.(Entry)

			// If the current time is after insertionTime + cleanupInterval, delete
			if time.Now().After(entryTyped.InsertionTime.Add(c.cleanupInterval)) {
				c.lruCache.Remove(k)
				log.Debug("Key has expired, removing", "k", k)
				continue
			}

			// finish when we've found a key that is not expired,
			// as the rest will be okay as well.
			return
		}
	}
}
