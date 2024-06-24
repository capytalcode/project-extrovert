package api

import (
	"io"
	"net/http"

	"extrovert/internals"
)

func RobotsTxt(w http.ResponseWriter, r *http.Request) {
	error := internals.HttpErrorHelper(w)

	aiList, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/robots.txt")
	if error("Error trying to fetch ai block list", err, http.StatusInternalServerError) {
		return
	}

	bytes, err := io.ReadAll(aiList.Body)
	if error("Error trying to read ai block list", err, http.StatusInternalServerError) {
		return
	}
	_, err = w.Write(bytes)
	if error("Error trying to write ai block list", err, http.StatusInternalServerError) {
		return
	}

	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")
	w.Header().Add("Content-Type", "text/plain")
}
