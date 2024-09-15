package storage

import (
	"sync"

	"github.com/timb418/url-shortener/internal/app/storageerrors"
)

type LinkStorage struct {
	originalToShort map[string]string
	shortToOriginal map[string]string
	mu              sync.RWMutex
}

func NewLinkStorage() *LinkStorage {
	return &LinkStorage{
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
	}
}

func (ls *LinkStorage) StoreLink(originalURL, shortURL string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Check if the original URL or short URL already exists
	if _, exists := ls.originalToShort[originalURL]; exists {
		return storageerrors.ErrOriginalUrlAlreadyExists
	}
	if _, exists := ls.shortToOriginal[shortURL]; exists {
		return storageerrors.ErrShortlUrlAlreadyExists
	}

	// Store the mappings
	ls.originalToShort[originalURL] = shortURL
	ls.shortToOriginal[shortURL] = originalURL

	return nil
}

func (ls *LinkStorage) GetOriginal(shortURL string) (string, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	originalURL, exists := ls.shortToOriginal[shortURL]
	if !exists {
		return "", storageerrors.ErrShortUrlNotFound
	}

	return originalURL, nil
}

func (ls *LinkStorage) GetShortened(originalURL string) (string, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	shortURL, exists := ls.originalToShort[originalURL]
	if !exists {
		return "", storageerrors.ErrOriginalUrlNotFound
	}

	return shortURL, nil
}
