package routes

import (
	e "errors"
	"io"
	"net/http"

	"extrovert/internals/router/errors"
)

type AITxt struct{}

func (_ AITxt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")

	list, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/ai.txt")
	if err != nil {
		errors.NewErrInternal(e.New("Unable to fetch ai.txt list"), err).ServeHTTP(w, r)
		return
	}

	bytes, err := io.ReadAll(list.Body)
	if err != nil {
		errors.NewErrInternal(e.New("Unable to read dynamic ai.txt list"), err).ServeHTTP(w, r)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	_, err = w.Write(bytes)
	if err != nil {
		errors.NewErrInternal(e.New("Unable to write ai.txt list"), err).ServeHTTP(w, r)
		return
	}
}

type RobotsTxt struct{}

func (_ RobotsTxt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")

	list, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/robots.txt")
	if err != nil {
		errors.NewErrInternal(e.New("Unable to fetch robots.txt list"), err).ServeHTTP(w, r)
		return
	}

	bytes, err := io.ReadAll(list.Body)
	if err != nil {
		errors.NewErrInternal(e.New("Unable to read dynamic robots.txt list"), err).ServeHTTP(w, r)
		return
	}

	_, err = io.WriteString(w, string(bytes))
	if err != nil {
		errors.NewErrInternal(e.New("Unable to write robots.txt list"), err).ServeHTTP(w, r)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
}
