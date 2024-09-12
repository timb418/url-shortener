package server

import "github.com/go-chi/chi"

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
