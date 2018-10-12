package proxy

import (
	"bytes"
	"net/http"

	"github.com/go-redis/redis"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

const (
	msgFailedHTTPShutdown = "Failed to shutdown http server"
)

// request represents some action that has hit our proxy
type request struct {
	Done       chan bool
	Key        string
	Resp       *bytes.Buffer
	StatusCode int
}

// NewRequest creates a new proxy request
func newRequest(key string, resp *bytes.Buffer) *request {
	return &request{
		Done:       make(chan bool),
		Key:        key,
		Resp:       resp,
		StatusCode: 0,
	}
}

// newWorker creates a new worker to listen on all proxy
// server requests coming into the system.
func (s *Server) newWorkerPool(count int) {
	for i := 0; i < count; i++ {
		go func() {
			for {
				select {
				case r := <-s.requests:
					log.Debug("Worker received request", "r", r)
					s.requestHandler(r)
					r.Done <- true

				case _ = <-s.done:
					return

				}
			}
		}()
	}

	log.Debug("All workers started!", "count", count)
}

func (s *Server) requestHandler(r *request) {
	r.StatusCode = http.StatusOK
	val, err := s.get(r.Key)
	if err != nil && err != redis.Nil {
		log.Error(msgProxiedGetFailed, "err", err)
		_, err := r.Resp.Write([]byte(msgProxiedGetFailed))
		if err != nil {
			log.Error(msgWriteFailed, "err", err)
		}
		r.StatusCode = http.StatusInternalServerError
		return
	}

	_, err = r.Resp.Write([]byte(val))
	if err != nil {
		log.Error(msgWriteFailed, "err", err)
		r.StatusCode = http.StatusInternalServerError
	}
}

// get returns the cached value if it exists,
// otherwise it will get the underlying storage value.
func (s *Server) get(key string) (string, error) {

	// Check local cache for key
	if val, found := s.Cache.Get(key); found {
		valStr, ok := val.(string)
		if !ok {
			return "", errors.New(msgCacheTypeAssertionFailed)
		}

		return valStr, nil
	}

	// Not in cache, get from redis and set to local cache
	val, err := s.StorageClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	s.Cache.Set(key, val)

	return val, nil
}
