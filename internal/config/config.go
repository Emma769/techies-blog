package config

import (
	"fmt"
	"strconv"
	"syscall"
)

type cfg struct{}

func New() *cfg {
	return &cfg{}
}

func (cfg cfg) LoadInt(key string, fallback int) int {
	value, ok := syscall.Getenv(key)

	if !ok {
		return fallback
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return i
}

func (cfg cfg) MustLoadString(key string) string {
	value, ok := syscall.Getenv(key)

	if !ok {
		panic(fmt.Sprintf("%s has no env value", key))
	}

	return value
}
