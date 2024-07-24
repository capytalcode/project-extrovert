package app

import (
	"extrovert/internals/oauth"
	"net/url"
	"os"
)

const DOMAIN = "http://localhost:7331/"

const TWITTER_REDIRECT = "/api/twitter/oauth2"

var TWITTER_APP = func() oauth.TwitterOAuth {
	ru, _ := url.Parse(DOMAIN)
	ru = ru.JoinPath(TWITTER_REDIRECT)

	c := oauth.NewTwitterOAuth(
		os.Getenv("TWITTER_CLIENT_ID"),
		os.Getenv("TWITTER_CLIENT_SECRET"),
		ru,
	)

	return c
}()
