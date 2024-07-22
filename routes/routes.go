package routes

import (
	"extrovert/internals/router"
)

var ROUTES = []router.Route{
	{Pattern: "/{$}", Handler: Homepage{}},
	{Pattern: "/robots.txt", Handler: RobotsTxt{}},
	{Pattern: "/ai.txt", Handler: AITxt{}},
}
