package middleware

import (
	"net/http"
	"strings"

	"github.com/emma769/techies-blog/internal/utils/funclib"
)

type CORSConfig struct {
	Origins,
	Headers,
	Methods []string
	AllowCredentials bool
}

func EnableCORS(config CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")

			origin := r.Header.Get("Origin")
			req := r.Header.Get("Access-Control-Request-Method")

			for i := range len(config.Origins) {
				if origin == config.Origins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if funclib.NonWhiteSpace(req) && r.Method == "OPTIONS" {
						methods := strings.Join(config.Methods, ", ")
						w.Header().Set("Access-Control-Allow-Methods", methods)

						headers := strings.Join(config.Headers, ", ")
						w.Header().Add("Access-Control-Allow-Headers", headers)

						w.WriteHeader(200)
						return
					}

					break
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
