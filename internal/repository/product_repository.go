package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, category_id, brand_id, name, description, price, stock FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.BrandID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, ID int) (*model.Product, error) {
	var p model.Product
	query := "SELECT id, category_id, brand_id, name, description, price, stock FROM products WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&p.ID, &p.CategoryID, &p.BrandID, &p.Name, &p.Description, &p.Price, &p.Stock)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Product.GetByID: %v", err)
	}

	return &p, nil
}
