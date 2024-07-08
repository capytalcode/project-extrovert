package routes

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"

	"extrovert/internals"
)

func AiTxt() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		aiList, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/ai.txt")
		if err != nil {
			return err
		}

		bytes, err := io.ReadAll(aiList.Body)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, string(bytes))
		return err
	})
}

func AiTxtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")

	error := internals.HttpErrorHelper(w)
	err := AiTxt().Render(context.Background(), w)
	if error("Error trying to create ai block list", err, http.StatusInternalServerError) {
		return
	}
	w.Header().Add("Content-Type", "text/plain")
}

func RobotsTxt() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		aiList, err := http.Get("https://raw.githubusercontent.com/ai-robots-txt/ai.robots.txt/main/robots.txt")
		if err != nil {
			return err
		}

		bytes, err := io.ReadAll(aiList.Body)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, string(bytes))
		return err
	})
}

func RobotsTxtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")

	error := internals.HttpErrorHelper(w)
	err := RobotsTxt().Render(context.Background(), w)
	if error("Error trying to create robots block list", err, http.StatusInternalServerError) {
		return
	}
	w.Header().Add("Content-Type", "text/plain")
}
