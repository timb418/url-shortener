package main

import (
	"log"
	"net/http"

	"github.com/timb418/url-shortener/internal/app/handlers"
	"github.com/timb418/url-shortener/internal/app/server"
)

func main() {
	srv := server.NewServer()
	handlers.RegisterRoutes(srv)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
