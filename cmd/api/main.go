package main

import (
	"log"

	"github.com/gonan98/ecom-pc-api/internal/config"
	"github.com/gonan98/ecom-pc-api/internal/database"
)

func main() {
	db, err := database.NewPostgres(config.Envs.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()

	server := NewServer(config.Envs.Port, db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
