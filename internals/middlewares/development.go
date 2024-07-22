package middlewares

import (
	"log"
	"net/http"
)

type DevelopmentMiddleware struct {
	Logger *log.Logger
}

func (m DevelopmentMiddleware) Serve(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.Logger.Printf("New request: %s", r.URL.Path)

		handler(w, r)

		w.Header().Del("Cache-Control")
		w.Header().Add("Cache-Control", "max-age=0")
	}
}
