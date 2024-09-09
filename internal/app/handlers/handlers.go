package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/timb418/url-shortener/internal/app/storage"
)

var linkStorage = storage.NewLinkStorage()

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", shortenGivenLink)
	mux.HandleFunc("/{id}", getFullLinkByShort)
}

func shortenGivenLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Получаем строку URL
	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	shortenedURL := generateShortLink(originalURL)
	err = linkStorage.StoreLink(originalURL, shortenedURL)
	if err != nil {
		log.Println("could not store link")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, shortenedURL)
}

func getFullLinkByShort(w http.ResponseWriter, r *http.Request) {
	shortId := r.PathValue("id")
	w.Header().Set("Content-Type", "text/plain")

	original, err := linkStorage.GetOriginal(shortId)

	if err != nil {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateShortLink(long string) string {
	hasher := sha256.New()
	hasher.Write([]byte(long))
	hash := hasher.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(hash)[:8]
	return shortURL
}
