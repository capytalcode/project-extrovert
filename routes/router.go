package routes

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
	{
		Pattern: "/api/twitter/oauth",
		Static:  false,
		Page:    TwitterOAuth(),
		Handler: TwitterOAuthHandler,
	},
	{
		Pattern: "/robots.txt",
		Static:  true,
		Page:    RobotsTxt(),
		Handler: RobotsTxtHandler,
	},
	{
		Pattern: "/ai.txt",
		Static:  true,
		Page:    AiTxt(),
		Handler: AiTxtHandler,
	},
}

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
