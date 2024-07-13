package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"extrovert/templates/pages"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/a-h/templ"
)

type Client interface {
	OAuthHandler(w http.ResponseWriter, r *http.Request)
}

type DefaultClient struct {
	name string
	tokenEndpoint url.URL
	id          string
	redirectUri string
}

func (c DefaultClient) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Fatalf("TODO-ERR missing code parameter")
	}

	req := c.tokenEndpoint

	q := req.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code_verifier", "challenge")
	q.Add("challenge_method", "plain")
	q.Add("code", code)
	q.Add("client_id", c.id)
	q.Add("redirect_uri", c.redirectUri)

	res, err := http.Post(req.String(), "application/x-www-form-urlencoded", bytes.NewReader([]byte("")))
	if err != nil {
		log.Fatalf("TODO-ERR trying to get token on %s, error:\n%s", req.Host, err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil || res.StatusCode != 200 {
		log.Fatalf("TODO-ERR trying to read body on %s, body:\n%s\n\nerror:\n%s", req.Host, body, err)
	}

	var token DefaultClientToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatalf("TODO-ERR trying to parse json body to token:\n%s", err)
	}

	cookie := http.Cookie{
		Name: strings.ToUpper("__Host-TOKEN-" + c.name),
		// Value:    token.String(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	err = pages.RedirectPopUp(
		"Logged into "+c.name+"!",
		"Your "+c.name+" account was succeffully logged into project-extrovert. ",
		templ.SafeURL("/index.html"),
	).Render(context.Background(), w)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("TODO-ERR trying to render static page:\n%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type DefaultClientToken struct {
	Type      string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
}
