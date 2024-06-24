package config

import (
	"net/http"

	"extrovert/api"
	"extrovert/internals"
	"extrovert/pages"
)

var ROUTES = []internals.Page{
	{Path: "index.html", Component: pages.Index()},
}

func APIROUTES(mux *http.ServeMux) {
	mux.HandleFunc("/robots.txt", api.RobotsTxt)
	mux.HandleFunc("/ai.txt", api.AiTxt)
}
