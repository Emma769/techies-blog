package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

type RecoverConfig struct {
	Logger *slog.Logger
}

func RecoverWithConfig(config RecoverConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					config.Logger.ErrorContext(
						r.Context(),
						"server error!",
						slog.String("detail", fmt.Sprintf("%v", err)),
					)

					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(500)
					w.Write([]byte(`{"detail": "internal server error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
