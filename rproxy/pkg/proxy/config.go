package proxy

import "time"

// Config represents a configuration for running our redis
// redis proxy server.
type Config struct {
	HTTPPort             string        `envconfig:"HTTP_PORT" default:"9090"`
	TCPPort              string        `envconfig:"TCP_PORT" default:"6379"`
	MaxConcurrentQueries int           `envconfig:"MAX_CONCURRENT_QUERIES" default:"500"`
	CacheExpiry          time.Duration `envconfig:"CACHE_EXPIRY" default:"5s"`
	CacheCleanup         time.Duration `envconfig:"CACHE_CLEANUP" default:"3s"`
	CacheCapacity        int           `envconfig:"CACHE_CAPACITY" default:"10"`
	StorageURI           string        `envconfig:"STORAGE_URI" default:"localhost:6379"`
	StoragePassword      string        `envconfig:"STORAGE_PASSWORD" default:""`
	WorkerCount          int           `envconfig:"WORKER_COUNT" default:"30"`
}
