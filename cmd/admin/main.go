package main

import (
	"context"
	"log"

	"github.com/gonan98/ecom-pc-api/internal/config"
	"github.com/gonan98/ecom-pc-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	pool, err := database.NewPostgres(config.Envs.DatabaseUrl)
	if err != nil {
		log.Fatal("Failed to connect to postgres")
	}
	defer pool.Close()

	ctx := context.Background()

	var roleID int
	if err := pool.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", "admin").Scan(&roleID); err != nil {
		log.Fatalf("Admin role does not exist in database, first run all migrations")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(config.Envs.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(ctx, "INSERT INTO users (first_name, last_name, email, password_hash, role_id) VALUES ($1, $2, $3, $4, $5)", "admin", "admin", config.Envs.AdminEmail, hash, roleID)
	if err != nil {
		log.Fatal("Could not create user admin")
	}

	log.Println("Admin created successfully")
}
