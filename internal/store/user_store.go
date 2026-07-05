package store

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user model.User, roleId int64) error {
	query := "INSERT INTO users (first_name, last_name, email, password_hash, role_id) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, roleId)
	return err
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := "SELECT id, email, password_hash, role_id FROM users WHERE email = $1"
	err := s.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.RoleID)

	if err != nil {
		return nil, fmt.Errorf("Error in User.GetByEmail")
	}

	return &user, nil
}

func (s *UserStore) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	var result bool
	err := s.db.QueryRow(ctx, query, email).Scan(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}
