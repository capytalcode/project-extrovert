package middlewares

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

type CookiesCryptoMiddleware struct {
	Key string
}

func (m CookiesCryptoMiddleware) Serve(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.decryptCookies(r)
		handler(w, r)
		m.encryptCookies(w)
	}
}

func (m CookiesCryptoMiddleware) encrypt(pt string) string {
	aes, err := aes.NewCipher([]byte(m.Key))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	ct := gcm.Seal(nonce, nonce, []byte(pt), nil)

	return base64.URLEncoding.EncodeToString(ct)
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

		v = m.encrypt(strings.Trim(v, "\""))

		c = cn + "=\"" + v + "\";" + attrs

		w.Header().Add("Set-Cookie", c)
	}
}

func (m CookiesCryptoMiddleware) decrypt(ct string) string {
	cb, err := base64.URLEncoding.DecodeString(ct)
	if err != nil {
		panic(err)
	}
	ct = string(cb)

	aes, err := aes.NewCipher([]byte(m.Key))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ct := ct[:nonceSize], ct[nonceSize:]

	pt, err := gcm.Open(nil, []byte(nonce), []byte(ct), nil)
	if err != nil {
		panic(err)
	}

	return string(pt)
}

func (m CookiesCryptoMiddleware) decryptCookies(r *http.Request) {
	rcookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, c := range rcookies {
		c.Value = m.decrypt(c.Value)
		r.AddCookie(c)
	}
}
