package storage

import "errors"

var (
	ErrorNotFound = errors.New("Not found")
	ErrorSave     = errors.New("Failed to save")
)
