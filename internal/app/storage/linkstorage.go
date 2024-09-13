package storage

import (
	"errors"
	"sync"
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
		return errors.New("original URL already exists")
	}
	if _, exists := ls.shortToOriginal[shortURL]; exists {
		return errors.New("short URL already exists")
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
		return "", errors.New("short URL not found")
	}

	return originalURL, nil
}

func (ls *LinkStorage) GetShortened(originalURL string) (string, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	shortURL, exists := ls.originalToShort[originalURL]
	if !exists {
		return "", errors.New("original URL not found")
	}

	return shortURL, nil
}
