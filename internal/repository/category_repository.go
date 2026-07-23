package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db DBTX
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]types.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, description FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]types.Category, 0)
	for rows.Next() {
		var cat types.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Exists(ctx context.Context, ID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)"
	result := false
	err := r.db.QueryRow(ctx, query, ID).Scan(&result)
	return result, err
}

func (r *CategoryRepository) GetByID(ctx context.Context, ID int) (*types.Category, error) {
	var category types.Category
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&category.ID, &category.Name, &category.Description)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Category.GetByID: %v", err)
	}

	return &category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, category *types.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2)"
	_, err := r.db.Exec(ctx, query, category.Name, category.Description)
	return err
}

func (r *CategoryRepository) Update(ctx context.Context, category *types.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	_, err := r.db.Exec(ctx, query, category.Name, category.Description, category.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, ID int) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := r.db.Exec(ctx, query, ID)
	return err
}
