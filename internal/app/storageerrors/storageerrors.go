package storageerrors

import "errors"

var (
	ErrOriginalURLAlreadyExists error = errors.New("original URL already exists")
	ErrShortlURLAlreadyExists error = errors.New("short URL already exists")
	ErrShortURLNotFound         error = errors.New("short URL not found")
	ErrOriginalURLNotFound      error = errors.New("original URL not found")
)
