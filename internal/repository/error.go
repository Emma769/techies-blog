package repository

import (
	"errors"
	"strings"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateKey   = errors.New("duplicate key")
	ErrUpdateConflict = errors.New("update conflict")
)

func DuplKey(s string) bool {
	return strings.Contains(s, "duplicate")
}
