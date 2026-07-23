package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/types"
)

type RoleRepository struct {
	db DBTX
}

func NewRoleRepository(db DBTX) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) GetByID(ctx context.Context, ID int) (*types.Role, error) {
	var role types.Role
	query := "SELECT id, name, description FROM roles WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		return nil, fmt.Errorf("role get by ID %d: %w", ID, err)
	}

	return &role, nil
}

func (r *RoleRepository) GetByName(ctx context.Context, roleName types.RoleName) (*types.Role, error) {
	var role types.Role
	err := r.db.QueryRow(ctx, "SELECT id, name, description FROM roles WHERE name = $1", roleName).Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		return nil, fmt.Errorf("role get by name %s: %w", roleName, err)
	}

	return &role, nil
}
