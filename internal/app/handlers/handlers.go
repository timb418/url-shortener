package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/timb418/url-shortener/internal/app/storage"
)

var linkStorage = storage.NewLinkStorage()
var hasher = sha256.New()

func RegisterRoutes(mux *chi.Mux, baseURL string) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ShortenGivenLink(w, r, baseURL)
	})
	mux.Get("/{id}", GetFullLinkByShort)
}

func ShortenGivenLink(w http.ResponseWriter, r *http.Request, baseURL string) {
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
	if originalURL == "" ||
		(!strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://")) {
		http.Error(w, "Wrong URL", http.StatusBadRequest)
		return
	}

	shortenedURL := generateShortLink(originalURL)
	err = linkStorage.StoreLink(originalURL, shortenedURL)
	if err != nil {
		log.Println("could not store link")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	fmt.Fprint(w, baseURL+"/"+shortenedURL)
}

func GetFullLinkByShort(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/")

	original, err := linkStorage.GetOriginal(shortID)

	if err != nil {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateShortLink(long string) string {
	hasher.Write([]byte(long))
	hash := hasher.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(hash)[:8]
	return shortURL
}
