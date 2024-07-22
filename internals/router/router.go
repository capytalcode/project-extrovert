package router

import (
	"log"
	"net/http"
	"strings"
)

type Route struct {
	Pattern  string
	Handler  http.Handler
	Children *[]Route
}

type Router struct {
	routes      []Route
	middlewares []Middleware
	mux         *http.ServeMux
	serveHTTP   http.HandlerFunc
}

func NewRouter(rs []Route) *Router {
	mux := http.NewServeMux()
	Router{}.registerAllRoutes("/", rs, mux)
	return &Router{rs, []Middleware{}, mux, mux.ServeHTTP}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.serveHTTP(w, req)
}

func (r *Router) AddMiddleware(m Middleware) {
	r.middlewares = append(r.middlewares, m)
	r.serveHTTP = r.wrapMiddleares(r.middlewares, r.serveHTTP)
}

func (router Router) wrapMiddleares(ms []Middleware, h http.HandlerFunc) http.HandlerFunc {
	fh := h.ServeHTTP
	for _, m := range ms {
		fh = m.Serve(fh)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		mw := NewMiddlewaredResponse(w)
		fh(mw, r)
		_, _ = mw.ReallyWriteHeader()
	}
}

func (router Router) registerAllRoutes(p string, rs []Route, mux *http.ServeMux) {
	for _, r := range rs {
		pattern := strings.Join([]string{
			strings.TrimSuffix(p, "/"),
			strings.TrimPrefix(r.Pattern, "/"),
		}, "/")
		log.Printf("registering route %s", pattern)

		mux.Handle(pattern, r.Handler)

		if r.Children != nil {
			router.registerAllRoutes(pattern, *r.Children, mux)
		}
	}
}
