package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/errors"
	"github.com/gonan98/ecom-pc-api/internal/util"
)

type contextKey string

const userContextKey contextKey = "user"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		claims, err := auth.ValidateJWT(token)

		if err != nil {
			util.WriteError(w, errors.NewAPIError(http.StatusUnauthorized, err))
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)

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

func GetUserClaims(ctx context.Context) (*auth.UserClaims, error) {
	claims, ok := ctx.Value(userContextKey).(*auth.UserClaims)
	if !ok {
		return nil, fmt.Errorf("User claims not found")
	}

	return claims, nil
}
