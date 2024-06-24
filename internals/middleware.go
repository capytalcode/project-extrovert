package internals

import (
	"log"
	"net/http"
)

type Middleware struct {
	handler http.Handler
	dev     bool
	noCache bool
	logger  *log.Logger
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.logger.Printf("Handling request. path=%s", r.URL.Path)

	if m.dev {
		r.URL.Scheme = "http"
	} else {
		r.URL.Scheme = "https"
	}

	m.handler.ServeHTTP(w, r)

	if m.noCache {
		w.Header().Del("Cache-Control")
		w.Header().Add("Cache-Control", "max-age=0")
	}
}

func NewMiddleware(handler http.Handler, dev bool, noCache bool, logger *log.Logger) *Middleware {
	return &Middleware{handler, dev, noCache, logger}
}
