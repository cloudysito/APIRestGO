package middleware

import (
	"net/http"
	"os"
)

func ValidarToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		tokenSecreto := os.Getenv("SECRET_TOKEN")

		if token != tokenSecreto {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
