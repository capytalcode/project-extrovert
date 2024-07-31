package middlewares

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

type CookiesCryptoMiddleware struct {
	key    string
	logger *log.Logger
}

func NewCookiesCryptoMiddleware(key string, logger *log.Logger) CookiesCryptoMiddleware {
	return CookiesCryptoMiddleware{key, logger}
}

func (m CookiesCryptoMiddleware) Serve(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.decryptCookies(r)
		handler(w, r)
		m.encryptCookies(w)
	}
}

func (m CookiesCryptoMiddleware) encrypt(pt string) (string, error) {
	aes, err := aes.NewCipher([]byte(m.key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	ct := gcm.Seal(nonce, nonce, []byte(pt), nil)

	return base64.URLEncoding.EncodeToString(ct), nil
}

func (m CookiesCryptoMiddleware) encryptCookies(w http.ResponseWriter) {
	cookies := w.Header().Values("Set-Cookie")
	w.Header().Del("Set-Cookie")
	for _, c := range cookies {
		cv := strings.Split(c, ";")

		var c, attrs string
		if len(cv) > 1 {
			c, attrs = cv[0], strings.Join(cv[1:], ";")
		} else {
			c = cv[0]
		}

		cn, v := strings.Split(c, "=")[0], strings.Split(c, "=")[1]

		v, err := m.encrypt(strings.Trim(v, "\""))
		if err != nil {
			m.logger.Panicf(
				"ERRO: Unable to encrypt cookie \"%s\", skipping. Error: %s",
				cn, err.Error(),
			)
			continue
		}

		c = cn + "=\"" + v + "\";" + attrs

		w.Header().Add("Set-Cookie", c)
	}
}

func (m CookiesCryptoMiddleware) decrypt(ct string) (string, error) {
	cb, err := base64.URLEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}
	ct = string(cb)

	aes, err := aes.NewCipher([]byte(m.key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ct := ct[:nonceSize], ct[nonceSize:]

	pt, err := gcm.Open(nil, []byte(nonce), []byte(ct), nil)
	if err != nil {
		return "", err
	}

	return string(pt), nil
}

func (m CookiesCryptoMiddleware) decryptCookies(r *http.Request) {
	rcookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, c := range rcookies {
		cv, err := m.decrypt(c.Value)
		if err != nil {
			m.logger.Panicf(
				"ERRO: Unable to decrypt cookie \"%s\", skipping: Error: %s",
				c.Name, err.Error(),
			)
			continue
		}
		c.Value = cv
		r.AddCookie(c)
	}
}
