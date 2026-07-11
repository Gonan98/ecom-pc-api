package config

import (
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	Port          string
	DatabaseUrl   string
	JWTSecret     string
	JWTExpiration string
	AdminEmail    string
	AdminPassword string
}

func initConfig() envConfig {
	godotenv.Load()

	return envConfig{
		Port:          getEnv("PORT", "8080"),
		DatabaseUrl:   getEnv("DATABASE_URL", "postgres://postgres:MyPostgrespassword@localhost:5432/ecomdb"),
		JWTSecret:     getEnv("JWT_SECRET", "mysecretpassword"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin12345"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

var Envs = initConfig()
