package internals

import (
	"net/http"
)

func HttpErrorHelper(w http.ResponseWriter) func(msg string, err error, status int) bool {
	return func(msg string, err error, status int) bool {
		if err != nil {
			w.WriteHeader(status)
			_, err = w.Write([]byte(msg + "\n Error: " + err.Error()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Error trying to return error code (somehow):\n" + err.Error()))
			}
			return true
		}
		return false
	}
}
