package proxy

import (
	"context"

	"net/http"

	"github.com/ggriffiths/gofun/rproxy/pkg/cache"
	"github.com/go-redis/redis"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

const (
	apiVersion    = "v1"
	redisEndpoint = "redis"

	msgProxyFailed              = "Proxy HTTP server failed"
	msgProxiedGetFailed         = "Proxy GET failed"
	msgCacheTypeAssertionFailed = "Type assertion for cache GET failed"
	msgMethodNotAllowed         = "This proxy only supports GET"
	msgWriteFailed              = "Proxy failed to write an HTTP response"
	msgNoKeyFound               = "You must pass a key to this endpoint"
)

// Proxy represents our proxy server
type Proxy struct {
	requests chan *request
	done     chan bool

	Server *http.Server
	Cache  *cache.Cache
	Client *redis.Client

	config *Config
}

// Start begins our proxy server
func Start(config *Config) error {
	// Create proxy
	proxy := Proxy{
		Cache: cache.New(config.CacheCapacity, config.CacheExpiry),
		Client: redis.NewClient(&redis.Options{
			Addr:     config.StorageURI,
			Password: config.StoragePassword,
			DB:       0, //default database
		}),
		config: config,
	}

	// Start proxy worker - listens on request events
	proxy.newWorkerPool(config.WorkerCount)

	// Create http server and start it up
	mux := http.NewServeMux()
	mux.HandleFunc(apiVersion+"/"+redisEndpoint, proxy.enqueue)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	return errors.Wrap(server.ListenAndServe(), msgProxyFailed)
}

// Close shuts down all request workers and the http server
func (p *Proxy) Close() {
	if err := p.Server.Shutdown(context.Background()); err != nil {
		log.Error(msgFailedHTTPShutdown, "err", err)
	}

	// Shutdown all workers
	for i := 0; i < p.config.WorkerCount; i++ {
		p.done <- true
	}
}

func (p *Proxy) enqueue(w http.ResponseWriter, r *http.Request) {
	p.requests <- newRequest(w, r)
}
