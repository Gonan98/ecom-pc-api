package config

import (
	"log"
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
		DatabaseURL:   getEnv("DATABASE_URL"),
		JWTSecret:     getEnv("JWT_SECRET"),
		JWTExpiration: getEnv("JWT_EXPIRATION"),
		AdminEmail:    getEnv("ADMIN_EMAIL"),
		AdminPassword: getEnv("ADMIN_PASSWORD"),
	}
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		log.Fatalf("env variable not found: %s", key)
	}

	return value
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

var Envs = initConfig()
