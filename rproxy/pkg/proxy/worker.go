package proxy

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	msgFailedHTTPShutdown = "Failed to shutdown http server"
)

// request represents some action that has hit our proxy
type request struct {
	httpWriter  http.ResponseWriter
	httpRequest *http.Request
}

// NewRequest creates a new proxy request
func newRequest(w http.ResponseWriter, r *http.Request) *request {
	return &request{
		httpWriter:  w,
		httpRequest: r,
	}
}

// newWorker creates a new worker to listen on all proxy
// server requests coming into the system.
func (p *Proxy) newWorkerPool(count int) {
	for id := 0; id < count; id++ {
		go func() {
			for {
				select {
				case r := <-p.requests:
					p.redisHandler(r.httpWriter, r.httpRequest)

				case _ = <-p.done:
					return
				}
			}
		}()
	}
}

func (p *Proxy) redisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.handleGet(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte(msgMethodNotAllowed)); err != nil {
			log.Error(msgMethodNotAllowed, "err", err)
		}

	}
}
