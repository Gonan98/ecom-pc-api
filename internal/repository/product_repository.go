package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) WithTx(tx pgx.Tx) *ProductRepository {
	return &ProductRepository{
		db: tx,
	}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]types.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, category_id, brand_id, name, description, image_url, price, stock, is_active FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.BrandID, &p.Name, &p.Description, &p.ImageUrl, &p.Price, &p.Stock, &p.IsActive); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, ID int) (*types.Product, error) {
	var p types.Product
	query := "SELECT id, category_id, brand_id, name, description, image_url, price, stock, is_active FROM products WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&p.ID, &p.CategoryID, &p.BrandID, &p.Name, &p.Description, &p.ImageUrl, &p.Price, &p.Stock, &p.IsActive)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Product.GetByID: %v", err)
	}

	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *types.Product) error {
	query := "INSERT INTO products (category_id, brand_id, name, description, image_url, price, stock) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, query, product.CategoryID, product.BrandID, product.Name, product.Description, product.ImageUrl, product.Price, product.Stock)
	return err
}

func (r *ProductRepository) Update(ctx context.Context, product *types.Product) error {
	q := "UPDATE products SET category_id = $1, brand_id = $2, name = $3, description = $4, image_url = $5, price = $6, stock = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP WHERE id = $9"
	_, err := r.db.Exec(ctx, q, product.CategoryID, product.BrandID, product.Name, product.Description, product.ImageUrl, product.Price, product.Stock, product.IsActive, product.ID)
	return err
}

func (r *ProductRepository) UpdateStock(ctx context.Context, productID int, stock int) error {
	q := "UPDATE products SET stock = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
	_, err := r.db.Exec(ctx, q, stock, productID)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, productID int) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := r.db.Exec(ctx, query, productID)
	return err
}
