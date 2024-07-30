package routes

import (
	"extrovert/internals/app"
	"extrovert/internals/router"
)

var ROUTES = []router.Route{
	{Pattern: "/{$}", Handler: Homepage{}},
	{Pattern: "/robots.txt", Handler: RobotsTxt{}},
	{Pattern: "/ai.txt", Handler: AITxt{}},
	{Pattern: app.TWITTER_REDIRECT, Handler: app.TWITTER_APP},
	{Pattern: app.MASTODON_REDIRECT, Handler: app.MASTODON_APP},
}
