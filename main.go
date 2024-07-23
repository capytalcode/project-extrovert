package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"extrovert/internals/middlewares"
	"extrovert/internals/router"
	"extrovert/routes"
)

var logger = log.Default()

func main() {
	port := flag.Int("p", 8080, "the port to run the server")
	dev := flag.Bool("d", false, "if the server is in development mode")

	flag.Parse()

	if *dev {
		log.Printf("Running server in DEVELOPMENT MODE")
	}

	r := router.NewRouter(routes.ROUTES)
	if *dev {
		r.AddMiddleware(middlewares.DevelopmentMiddleware{Logger: logger})
	}
	r.AddMiddleware(middlewares.CookiesCryptoMiddleware{os.Getenv("CRYPTO_COOKIES_KEY")})

	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), r)
	if err != nil {
		logger.Fatalf("Server crashed due to:\n%s", err)
	}
}
