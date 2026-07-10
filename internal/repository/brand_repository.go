package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/jackc/pgx/v5"
)

type BrandRepository struct {
	db DBTX
}

func NewBrandRepository(db DBTX) *BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

func (r *BrandRepository) GetAll(ctx context.Context) ([]types.Brand, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, website FROM brands ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brands := make([]types.Brand, 0)
	for rows.Next() {
		var b types.Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Website); err != nil {
			return nil, err
		}
		brands = append(brands, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *BrandRepository) GetByID(ctx context.Context, ID int) (*types.Brand, error) {
	var brand types.Brand
	query := "SELECT id, name, website FROM brands WHERE id = $1"
	err := r.db.QueryRow(ctx, query, ID).Scan(&brand.ID, &brand.Name, &brand.Website)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Brand.GetByID: %v", err)
	}

	return &brand, nil
}
