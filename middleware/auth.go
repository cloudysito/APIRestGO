package middleware

import (
	"net/http"
	"os"
)

func ValidateToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		secretToken := os.Getenv("SECRET_TOKEN")

		if token != secretToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
