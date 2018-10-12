package proxy

import (
	"bytes"
	"context"
	"net"
	"os"
	"strings"

	"net/http"

	"github.com/ggriffiths/gofun/rproxy/pkg/cache"
	"github.com/go-redis/redis"
	log "github.com/inconshreveable/log15"
)

const (
	apiVersion    = "v1"
	redisEndpoint = "redis"

	msgFailedCacheCreate        = "Failed to create cache"
	msgProxyHTTPFailed          = "Proxy HTTP server failed"
	msgProxyTCPFailed           = "Proxy TCP server failed"
	msgProxiedGetFailed         = "Proxy GET failed"
	msgCacheTypeAssertionFailed = "Type assertion for cache GET failed"
	msgMethodNotAllowed         = "This proxy only supports GET"
	msgInvalidRedisRequest      = "Invalid TCP redis request"
	msgWriteFailed              = "Proxy failed to write an HTTP response"
	msgNoKeyFound               = "You must pass a key to this endpoint"
)

// Server represents our proxy server
type Server struct {
	requests chan *request
	done     chan bool

	HTTPServer    *http.Server
	Cache         *cache.Cache
	StorageClient *redis.Client

	config *Config
}

// NewServer creates a new proxy server
func NewServer(config *Config) *Server {
	cache, err := cache.New(config.CacheCapacity, config.CacheExpiry, config.CacheCleanup)
	if err != nil {
		log.Crit(msgFailedCacheCreate, "err", err)
		os.Exit(1)
	}

	server := &Server{
		requests: make(chan *request, config.MaxConcurrentQueries),
		done:     make(chan bool),

		Cache: cache,
		StorageClient: redis.NewClient(&redis.Options{
			Addr:     config.StorageURI,
			Password: config.StoragePassword,
			DB:       0, //default database
		}),

		config: config,
	}

	return server
}

// Start begins our proxy server
func (s *Server) Start() error {
	// Close when start finishes
	defer s.Close()

	// Start worker pool
	s.newWorkerPool(s.config.WorkerCount)

	// Start HTTP server in separate goroutine
	go s.startHTTPServer()

	// Start TCP server in main goroutine
	return s.startTCPServer()
}

func (s *Server) startHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/"+apiVersion+"/"+redisEndpoint, s.handleHTTPRequest)
	httpServer := &http.Server{
		Addr:    ":" + s.config.HTTPPort,
		Handler: mux,
	}
	s.HTTPServer = httpServer

	log.Debug("Starting HTTP server")
	err := s.HTTPServer.ListenAndServe()
	if err != nil {
		log.Error(msgProxyHTTPFailed, "err", err)
		os.Exit(1)
	}
}

func (s *Server) startTCPServer() error {
	l, err := net.Listen("tcp", "rproxy:"+s.config.TCPPort)
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		log.Debug("Starting TCP server")
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		// Handle connections in a new goroutine.
		go s.handleTCPRedisRequest(conn)
	}
}

func (s *Server) handleTCPRedisRequest(conn net.Conn) {
	defer conn.Close()

	// Make a buffer to hold incoming data.
	log.Debug("Handling TCP request...")
	req := make([]byte, 1024)

	// Read tcp request into req
	reqLen, err := conn.Read(req)
	if err != nil {
		log.Error("Failed to read redis request", "err", err)
		return
	}

	// break up req into parts, handle GET
	parts := strings.Split(string(req[:reqLen]), "\r\n")
	if len(parts) < 4 {
		conn.Write([]byte(msgInvalidRedisRequest))
	}
	switch parts[2] {
	case "get":
		key := parts[4]

		var resp = bytes.NewBuffer([]byte{})
		proxyRequest := newRequest(key, resp)
		s.requests <- proxyRequest

		// Block until request finished
		<-proxyRequest.Done

		// We'll communicate internally with http status code for now.
		if proxyRequest.StatusCode != http.StatusOK {
			redisError := "-" + string(resp.Bytes()) + "\r\n"
			conn.Write([]byte(redisError))
		}

		// Write response when done
		redisResponse := ":" + string(resp.Bytes()) + "\r\n"
		_, err = conn.Write([]byte(redisResponse))
		if err != nil {
			log.Error("Failed to write redis response", "err", err)
			return
		}

	default:
		conn.Write([]byte(msgMethodNotAllowed))
		return
	}
}

func (s *Server) handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	if len(s.requests) >= s.config.MaxConcurrentQueries {
		w.Write([]byte("Max concurrent requests hit, retry later"))
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if key == "" {
			w.Write([]byte("Key query param must be nonempty"))
			w.WriteHeader(http.StatusBadRequest)
		}

		resp := bytes.NewBuffer([]byte{})
		proxyRequest := newRequest(key, resp)
		s.requests <- proxyRequest

		// Block until finished writing http.ResponseWriter inside
		// of our worker.
		<-proxyRequest.Done
		w.Write(proxyRequest.Resp.Bytes())
		w.WriteHeader(proxyRequest.StatusCode)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte(msgMethodNotAllowed)); err != nil {
			log.Error(msgMethodNotAllowed, "err", err)
		}

	}

}

// Close shuts down all request workers and the http server
func (s *Server) Close() {
	// Shutdown HTTP server
	if err := s.HTTPServer.Shutdown(context.Background()); err != nil {
		log.Error(msgFailedHTTPShutdown, "err", err)
	}

	// Close cache
	s.Cache.Close()

	// Shutdown all workers
	for i := 0; i < s.config.WorkerCount; i++ {
		s.done <- true
	}
}
