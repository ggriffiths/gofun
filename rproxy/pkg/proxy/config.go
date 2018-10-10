package proxy

import "time"

// Config represents a configuration for running our redis
// redis proxy server.
type Config struct {
	Port            string
	CacheExpiry     time.Duration
	CacheCapacity   int64
	StorageURI      string
	StoragePassword string
	WorkerCount     int
}
