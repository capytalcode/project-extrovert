package pages

import (
	"net/http"

	"github.com/a-h/templ"
)

var ROUTES = []Route{
	{
		Pattern: "/index.html",
		Static:  true,
		Page:    IndexPage(),
		Handler: IndexHandler,
	},

type RouteHandler = func(http.ResponseWriter, *http.Request)

type Route struct {
	Pattern string
	Static  bool
	Handler RouteHandler
	Page    templ.Component
}

func RegisterAllRoutes(routes []Route, s *http.ServeMux) {
	for _, r := range routes {
		s.HandleFunc(r.Pattern, r.Handler)
	}
}
