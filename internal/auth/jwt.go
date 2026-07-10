package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gonan98/ecom-pc-api/internal/config"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (u *UserClaims) UserID() (int, error) {
	return strconv.Atoi(u.Subject)
}

func GenerateJWT(userID int, role string) (string, error) {
	exp, err := time.ParseDuration(config.Envs.JWTExpiration)
	if err != nil {
		return "", err
	}

	claims := UserClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateJWT(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.Envs.JWTSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, err
}

func GetClaims(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(UserContextKey).(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("User claims not found")
	}

	return claims, nil
}
