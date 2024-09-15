package storageerrors

import "errors"

var (
	ErrOriginalUrlAlreadyExists error = errors.New("original URL already exists")
	ErrShortlUrlAlreadyExists error = errors.New("short URL already exists")
	ErrShortUrlNotFound         error = errors.New("short URL not found")
	ErrOriginalUrlNotFound      error = errors.New("original URL not found")
)
