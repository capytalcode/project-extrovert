package main

import (
	"context"
	"flag"
	"log"

	"extrovert/config"
	"extrovert/internals"
)

func main() {
	dir := flag.String("d", "./dist", "the directory to write the files")
	staticDir := flag.String("s", "./static", "the directory to copy static files from")

	flag.Parse()

	w := internals.StaticWriter{
		DistDir:   dir,
		StaticDir: staticDir,
		Pages:     config.ROUTES,
		Context:   context.Background(),
		Logger:    *log.Default(),
	}

	err := w.WriteAll()
	if err != nil {
		log.Fatal(err)
	}
}
