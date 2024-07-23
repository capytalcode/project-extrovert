package middlewares

import (
	// "extrovert/internals/cookies"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

const KEY = "AGVSBG93B3JSZAAA"

const STRING = "Hello world, I'm a Cookie"

func handler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("__Host-test")

	if c.Value != STRING {
		log.Fatalf(
			"Request cookie wasn't correctly decrypted.\nOriginal: %s\nDecrypted: %s",
			STRING, c.Value,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-response-test",
		Value:    STRING,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   2000,
		Domain:   "localhost",
		HttpOnly: true,
	})
}

func TestRequest(t *testing.T) {
	m := NewCookiesCryptoMiddleware(KEY, log.Default())

	e, err := m.encrypt(STRING)
	if err != nil {
		panic(err)
	}

	var COOKIE = http.Cookie{
		Name:     "__Host-test",
		Value:    e,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   2000,
		Domain:   "localhost",
		HttpOnly: true,
	}

	req := httptest.NewRequest("GET", "https://localhost:3030", nil)
	req.AddCookie(&COOKIE)

	w := httptest.NewRecorder()

	h := m.Serve(handler)
	h(w, req)
	w.WriteHeader(http.StatusOK)
	log.Printf("%#v", w.Header().Values("Set-Cookie"))

	res := w.Result()
	cs := res.Cookies()
	log.Print(cs)
	ci := slices.IndexFunc(cs, func(c *http.Cookie) bool {
		return c.Name == "__Host-response-test"
	})
	c := cs[ci]

	log.Print(c)

	d, err := m.decrypt(c.Value)
	if err != nil {
		panic(err)
	}

	if d != STRING {
		log.Fatalf(
			"Response cookie wasn't correctly encrypted.\nOriginal: %s\nEncrypted: %s",
			STRING, d,
		)
	}
}

func TestEncrypt(t *testing.T) {
	m := NewCookiesCryptoMiddleware(KEY, log.Default())

	e, err := m.encrypt(STRING)
	if err != nil {
		panic(err)
	}

	log.Printf("Encrypted %s", e)

	d, err := m.decrypt(e)
	if err != nil {
		panic(err)
	}

	log.Printf("Decrypted %s", d)

	if d != STRING {
		log.Fatalf("Decrypted value isn't equal to original.\n"+
			"Original: %s\n"+
			"Decrypted: %s\n"+
			"\n"+
			"Encrypted: %s\n",
			STRING, d, e)
	}
}
