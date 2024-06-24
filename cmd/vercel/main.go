package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"extrovert/config"
	"extrovert/internals"
)

type VercelConfig struct {
	OutputDirectory string `json:"outputDirectory"`
}

var logger = log.Default()

func main() {
	configPath := flag.String("c", "./vercel.json", "the path to the vercel.json file")
	staticDir := flag.String("s", "./static", "the directory to copy static files from")
	port := flag.Int("p", 8080, "the port to run the server")

	flag.Parse()

	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		logger.Fatalf("Unable to read vercel.json file due to:\n%s", err)
	}

	var c VercelConfig
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		logger.Fatalf("Unable to parse vercel.json file due to:\n%s", err)
	}

	w := internals.StaticWriter{
		DistDir:   &c.OutputDirectory,
		StaticDir: staticDir,
		Pages:     config.ROUTES,
		Context:   context.Background(),
		Logger:    *log.Default(),
	}

	logger.Print("Writing static files")
	err = w.WriteAll()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Print("Starting server")
	mux := http.NewServeMux()

	config.APIROUTES(mux)
	mux.Handle("/", http.FileServer(http.Dir(c.OutputDirectory)))

	logger.Printf("Running server at port: %v", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), mux)
	if err != nil {
		logger.Fatalf("Server crashed due to:\n%s", err)
	}
}
