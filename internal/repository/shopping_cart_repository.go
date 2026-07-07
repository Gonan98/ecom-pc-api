package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShoppingCartRepository struct {
	db *pgxpool.Pool
}

func NewShoppingCartRepository(db *pgxpool.Pool) *ShoppingCartRepository {
	return &ShoppingCartRepository{
		db: db,
	}
}

func (r *ShoppingCartRepository) Create(ctx context.Context, userID int) (int, error) {
	var cartID int
	query := "INSERT INTO shopping_carts (user_id) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(ctx, query, userID).Scan(&cartID)

	if err != nil && err != pgx.ErrNoRows {
		return 0, fmt.Errorf("ShoppingCartRepository.Create: %v", err)
	}

	return cartID, err
}

func (r *ShoppingCartRepository) Get(ctx context.Context, userID int) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	query := "SELECT id, user_id FROM shopping_carts WHERE user_id = $1"
	err := r.db.QueryRow(ctx, query, userID).Scan(&cart.ID, &cart.UserID)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("ShoppingCartRepository.Get: %v", err)
	}

	return &cart, nil
}
