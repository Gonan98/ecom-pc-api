package middleware

import (
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.GetClaims(r.Context())
		if err != nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		if claims.Role != types.RoleNameAdmin {
			http.Error(w, "Not allowed, you are not an administrator", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
