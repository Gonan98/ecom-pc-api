package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gonan98/ecom-pc-api/internal/config"
)

func GenerateJWT(userID int, role string) (string, error) {
	exp, err := time.ParseDuration(config.Envs.JWTExpiration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  strconv.Itoa(userID),
		"role": role,
		"exp":  time.Now().Add(exp).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

// func GetClaimsFromContext(ctx context.Context) (int, string) {
// 	// Get userId y role

// }
