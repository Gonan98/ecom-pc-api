package middleware

import (
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/auth"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.GetClaims(r.Context())
		if err != nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, "Not allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
