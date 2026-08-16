package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

func BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedToken := os.Getenv("API_TOKEN")

		if expectedToken == "" {
			http.Error(
				w,
				"server authentication is not configured",
				http.StatusInternalServerError,
			)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"missing authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(
				w,
				"invalid authorization scheme",
				http.StatusUnauthorized,
			)
			return
		}

		providedToken := strings.TrimSpace(
			strings.TrimPrefix(authHeader, prefix),
		)

		if providedToken == "" {
			http.Error(
				w,
				"missing bearer token",
				http.StatusUnauthorized,
			)
			return
		}

		if subtle.ConstantTimeCompare(
			[]byte(providedToken),
			[]byte(expectedToken),
		) != 1 {

			http.Error(
				w,
				"invalid token",
				http.StatusUnauthorized,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}