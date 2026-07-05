package store

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleStore struct {
	db *pgxpool.Pool
}

func NewRoleStore(db *pgxpool.Pool) *RoleStore {
	return &RoleStore{
		db: db,
	}
}

func (s *RoleStore) GetByID(ctx context.Context, ID int) (*model.Role, error) {
	var role model.Role
	query := "SELECT id, name FROM roles WHERE id = $1"
	err := s.db.QueryRow(ctx, query, ID).Scan(&role.ID, &role.Name)

	if err != nil {
		return nil, fmt.Errorf("Error in Role.GetByID")
	}

	return &role, nil
}

func (s *RoleStore) GetCustomerRoleID(ctx context.Context) (int64, error) {
	var customerId int64
	err := s.db.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", "customer").Scan(&customerId)

	if err != nil {
		return 0, err
	}

	return customerId, nil
}
