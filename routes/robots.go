package routes

import (
	"io"
	"net/http"

	"extrovert/internals"
)

type AITxt struct{}

func (_ AITxt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")

	error := internals.HttpErrorHelper(w)

	aiList, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/ai.txt")
	if error("Error trying to create ai block list", err, http.StatusInternalServerError) {
		return
	}

	bytes, err := io.ReadAll(aiList.Body)
	if error("Error trying to create ai block list", err, http.StatusInternalServerError) {
		return
	}
	_, err = io.WriteString(w, string(bytes))
	if error("Error trying to create ai block list", err, http.StatusInternalServerError) {
		return
	}

	w.Header().Add("Content-Type", "text/plain")
}

type RobotsTxt struct{}

func (_ RobotsTxt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")

	error := internals.HttpErrorHelper(w)
	aiList, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/robots.txt")
	if error("Error trying to create robots block list", err, http.StatusInternalServerError) {
		return
	}

	bytes, err := io.ReadAll(aiList.Body)
	if error("Error trying to create robots block list", err, http.StatusInternalServerError) {
		return
	}

	_, err = io.WriteString(w, string(bytes))
	if error("Error trying to create robots block list", err, http.StatusInternalServerError) {
		return
	}
	w.Header().Add("Content-Type", "text/plain")
}
