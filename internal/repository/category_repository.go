package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, description FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var cat model.Category
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

func (r *CategoryRepository) GetByID(ctx context.Context, ID int) (*model.Category, error) {
	var category model.Category
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&category.ID, &category.Name, &category.Description)

	if err != nil {
		return nil, fmt.Errorf("Category.GetByID: %v", err)
	}

	return &category, nil
}
