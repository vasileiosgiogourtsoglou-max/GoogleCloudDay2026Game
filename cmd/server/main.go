package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed public/*
var publicFS embed.FS

func main() {
	staticFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatalf("server: sub fs error: %v", err)
	}

	http.Handle("/", http.FileServer(http.FS(staticFS)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Olympus Game Server listening on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
