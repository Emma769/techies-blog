package services

import "errors"

var (
	ErrDuplicateKey   = errors.New("duplicate key")
	ErrUpdateConflict = errors.New("update conflict")
	ErrNotFound       = errors.New("not found")
)
