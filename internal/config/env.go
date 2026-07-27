package config

import (
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration string
	AdminEmail    string
	AdminPassword string
}

func initConfig() envConfig {
	godotenv.Load()

	return envConfig{
		Port:          getEnvOrDefault("PORT", "8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

var Envs = initConfig()
