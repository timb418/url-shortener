package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/timb418/url-shortener/internal/app/handlers"
	"github.com/timb418/url-shortener/internal/app/server"
)

func main() {
	router := server.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	handlers.RegisterRoutes(router)

	log.Println("Server initialized")
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
