package proxy

import (
	"net/http"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

// handleGet handles a single http.Get request to the proxy
func (p *Proxy) handleGet(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok {
		_, err := w.Write([]byte(msgNoKeyFound))
		if err != nil {
			log.Error(msgWriteFailed, "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
	}

	val, err := p.get(keys[0])
	if err != nil {
		log.Error(msgProxiedGetFailed, "err", err)
		_, err := w.Write([]byte(msgProxiedGetFailed))
		if err != nil {
			log.Error(msgWriteFailed, "err", err)
		}
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte(val))
	if err != nil {
		log.Error(msgWriteFailed, "err", err)
	}
}

// get returns the cached value if it exists,
// otherwise it will get the underlying storage value.
func (p *Proxy) get(key string) (string, error) {
	if val, found := p.Cache.Get(key); found {
		valStr, ok := val.(string)
		if !ok {
			return "", errors.New(msgCacheTypeAssertionFailed)
		}

		return valStr, nil
	}

	return p.Client.Get(key).Result()
}
