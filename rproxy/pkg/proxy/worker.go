package proxy

import (
	"net/http"

	"github.com/nu7hatch/gouuid"

	log "github.com/inconshreveable/log15"
)

const (
	msgFailedHTTPShutdown = "Failed to shutdown http server"
)

// request represents some action that has hit our proxy
type request struct {
	correlationID string
	httpWriter    http.ResponseWriter
	httpRequest   *http.Request
}

// NewRequest creates a new proxy request
func newRequest(w http.ResponseWriter, r *http.Request) *request {
	var correlationID string

	id, err := uuid.NewV4()
	if err != nil {
		// Should almost never fail, but if it does,
		// continuing operating though logging as an error
		log.Error("Failed to generate UUID", "err", err)
		correlationID = ""
		return &request{
			correlationID: "failed",
			httpWriter:    w,
			httpRequest:   r,
		}
	}
	correlationID = id.String()

	return &request{
		correlationID: correlationID,
		httpWriter:    w,
		httpRequest:   r,
	}
}

// newWorker creates a new worker to listen on all proxy
// server requests coming into the system.
func (p *Proxy) newWorkerPool(count int) {
	for i := 0; i < count; i++ {
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
