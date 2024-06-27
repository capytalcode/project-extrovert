package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"extrovert/internals"
	"extrovert/routes"
)

var logger = log.Default()

func main() {
	staticDir := flag.String("s", "./static", "the directory to copy static files from")
	port := flag.Int("p", 8080, "the port to run the server")
	dev := flag.Bool("d", false, "if the server is in development mode")
	cache := flag.Bool("c", true, "if the static files are cached")

	flag.Parse()

	if *dev {
		log.Printf("Running server in DEVELOPMENT MODE")
	}

	mux := http.NewServeMux()

	routes.RegisterAllRoutes(routes.ROUTES, mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			logger.Printf("Handling file server request. path=%s", r.URL.Path)
			http.FileServer(http.Dir(*staticDir)).ServeHTTP(w, r)
			return
		}
	})

	logger.Printf("Running server at port: %v", *port)

	middleware := internals.NewMiddleware(mux, *dev, !*cache, log.Default())
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), middleware)
	if err != nil {
		logger.Fatalf("Server crashed due to:\n%s", err)
	}
}
