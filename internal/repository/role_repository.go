package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) GetByID(ctx context.Context, ID int) (*model.Role, error) {
	var role model.Role
	query := "SELECT id, name, description FROM roles WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		return nil, fmt.Errorf("Role.GetByID: %v", err)
	}

	return &role, nil
}

func (r *RoleRepository) GetCustomerRoleID(ctx context.Context) (int, error) {
	customerID := 0
	err := r.db.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", "customer").Scan(&customerID)
	return customerID, err
}
