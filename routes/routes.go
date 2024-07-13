package routes

import (
	"context"
	"extrovert/internals"
	"extrovert/templates/pages"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func NewStaticPageHandler(c templ.Component) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.Render(context.Background(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("TODO-ERR trying to render static page:\n%s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

var ROUTES = []internals.Route{
	{
		Pattern: "/index.html",
		Static:  true,
		Page:    pages.Homepage(),
		Handler: NewStaticPageHandler(pages.Homepage()),
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
