package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user model.User, roleId int64) error {
	query := "INSERT INTO users (first_name, last_name, email, password_hash, role_id) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(ctx, query, user.FirstName, user.LastName, user.Email, user.PasswordHash, roleId)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, ID int) (*model.User, error) {
	var user model.User
	query := "SELECT id, first_name, last_name, email, role_id FROM users WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.RoleID)

	if err != nil {
		return nil, fmt.Errorf("UserStore.GetByEmail: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := "SELECT id, email, password_hash, role_id FROM users WHERE email = $1"
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.RoleID)

	if err != nil {
		return nil, fmt.Errorf("UserStore.GetByEmail: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	var result bool
	err := r.db.QueryRow(ctx, query, email).Scan(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}
