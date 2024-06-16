package middleware

import (
	"log/slog"
	"net/http"
)

type LoggerConfig struct {
	Logger *slog.Logger
}

func LoggerWithConfig(config LoggerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			config.Logger.InfoContext(
				r.Context(),
				"incoming request",
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
			)

			next.ServeHTTP(w, r)
		})
	}
}
