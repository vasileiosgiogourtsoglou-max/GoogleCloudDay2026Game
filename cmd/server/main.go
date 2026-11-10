package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dir := "public"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = "."
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Printf("Olympus Game Server listening on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
