package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/timb418/url-shortener/internal/app/config"
	"github.com/timb418/url-shortener/internal/app/handlers"
	"github.com/timb418/url-shortener/internal/app/server"
)

func main() {
	cfg := config.NewConfig()

	router := server.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	handlers.RegisterRoutes(router, cfg.BaseURL)

	log.Printf("Starting server on %s\n", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
