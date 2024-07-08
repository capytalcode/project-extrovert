package internals

import (
	"fmt"
	"net/http"
	"slices"
)

func RemoveDuplicates[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetCookie(name string, w http.ResponseWriter, r *http.Request) *http.Cookie {
	name = fmt.Sprintf("__Host-%s-%s-%s", APP_NAME, APP_VERSION, name)

	c := r.Cookies()
	i := slices.IndexFunc(c, func(c *http.Cookie) bool {
		return c.Name == name
	})
	var cookie *http.Cookie
	if i == -1 {
		cookie = &http.Cookie{
			Name:     name,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Secure:   true,
		}
	} else {
		cookie = c[i]
	}
	return cookie
}

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
