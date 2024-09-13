package server

import (
	"net/http"
)

func NewServer() *http.ServeMux {
	return http.NewServeMux()
}
