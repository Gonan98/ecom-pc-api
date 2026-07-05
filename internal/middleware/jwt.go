package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/errors"
	"github.com/gonan98/ecom-pc-api/internal/util"
)

func JWTMiddleware(apiHandler util.APIHandler) util.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenString := getTokenFromRequest(r)
		token, err := auth.ValidateJWT(tokenString)

		if err != nil {
			return errors.NewAPIError(http.StatusUnauthorized, fmt.Errorf("No token provided"))
		}

		if !token.Valid {
			return errors.NewAPIError(http.StatusForbidden, fmt.Errorf("Invalid token"))
		}

		claims := token.Claims.(jwt.MapClaims)
		sub := claims["sub"].(string)
		userID, err := strconv.Atoi(sub)
		if err != nil {
			return err
		}
		role := claims["role"].(string)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "role", role)
		r = r.WithContext(ctx)

		return apiHandler(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if token != "" {
		return token
	}

	return ""
}
