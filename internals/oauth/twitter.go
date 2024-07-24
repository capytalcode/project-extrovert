package oauth

import "net/url"

type TwitterOAuth struct {
	DefaultOAuthClient
}

func NewTwitterOAuth(id string, secret string, redirect *url.URL) TwitterOAuth {
	u, _ := url.Parse("https://api.twitter.com/2/")
	c := NewDefaultOAuthClient(u, id, secret, redirect)

	c.Name = "Twitter"

	u, _ = url.Parse("https://twitter.com/i/oauth2/authorize")
	q := c.AuthEndpoint.Query()
	q.Del("scope")
	q.Add("scope", "tweet.read tweet.write users.read")
	u.RawQuery = q.Encode()
	c.AuthEndpoint = u

	return TwitterOAuth{c}
}
