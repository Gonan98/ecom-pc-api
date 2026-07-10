package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gonan98/ecom-pc-api/internal/auth"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		claims, err := auth.ValidateJWT(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getToken(r *http.Request) string {
	header := r.Header.Get("Authorization")

	if header == "" {
		return ""
	}

	bearerToken := strings.Split(header, " ")
	if len(bearerToken) != 2 {
		return ""
	}

	if bearerToken[0] != "Bearer" {
		return ""
	}

	return bearerToken[1]
}
