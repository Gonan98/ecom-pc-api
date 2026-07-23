package repository

import (
	"context"
	"fmt"

	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/jackc/pgx/v5"
)

type CartRepository struct {
	db DBTX
}

func NewCartRepository(db DBTX) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) WithTx(tx pgx.Tx) *CartRepository {
	return &CartRepository{
		db: tx,
	}
}

func (r *CartRepository) Create(ctx context.Context, userID int) error {
	q := "INSERT INTO shopping_carts (user_id) VALUES ($1)"
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

func (r *CartRepository) CreateItem(ctx context.Context, item *types.CartItem) error {
	q := "INSERT INTO shopping_cart_items VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, q, item.CartID, item.ProductID, item.Quantity)
	return err
}

func (r *CartRepository) GetByUser(ctx context.Context, userID int) (*types.Cart, error) {
	var cart types.Cart
	q := "SELECT id, user_id FROM shopping_carts WHERE user_id = $1"
	err := r.db.QueryRow(ctx, q, userID).Scan(&cart.ID, &cart.UserID)

	if err != nil {
		return nil, fmt.Errorf("cart get by user %d: %w", userID, err)
	}

	return &cart, nil
}

func (r *CartRepository) GetItemsByUser(ctx context.Context, userID int) ([]types.CartItem, error) {
	q := "SELECT ci.cart_id, ci.product_id, ci.quantity FROM shopping_cart_items ci JOIN shopping_carts c ON ci.cart_id = c.id WHERE c.user_id = $1"
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]types.CartItem, 0)
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(&item.CartID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *CartRepository) ExistsItemsInCart(ctx context.Context, cartID int) (bool, error) {
	result := false
	q := "SELECT EXISTS (SELECT 1 FROM shopping_cart_items WHERE cart_id = $1)"
	err := r.db.QueryRow(ctx, q, cartID).Scan(&result)
	return result, err
}

func (r *CartRepository) ExistsItemInCartByProductID(ctx context.Context, cartID int, productID int) (bool, error) {
	result := false
	q := "SELECT EXISTS (SELECT 1 FROM shopping_cart_items WHERE cart_id = $1 AND product_id = $2)"
	err := r.db.QueryRow(ctx, q, cartID, productID).Scan(&result)
	return result, err
}

func (r *CartRepository) UpdateItemQuantity(ctx context.Context, cartID int, productID int, quantity int) error {
	q := "UPDATE shopping_cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3"
	_, err := r.db.Exec(ctx, q, quantity, cartID, productID)
	return err
}

func (r *CartRepository) DeleteItems(ctx context.Context, cartID int) error {
	q := "DELETE FROM shopping_cart_items WHERE cart_id = $1"
	_, err := r.db.Exec(ctx, q, cartID)
	return err
}

func (r *CartRepository) DeleteItemsByProductID(ctx context.Context, cartID int, productID int) error {
	q := "DELETE FROM shopping_cart_items WHERE cart_id = $1 AND product_id = $2"
	_, err := r.db.Exec(ctx, q, cartID, productID)
	return err
}

// query := SELECT
// 			c.product_id,
// 			p.name,
// 			p.price,
// 			ci.quantity
// 		FROM shopping_carts c
// 		JOIN shopping_cart_items ci ON ci.cart_id = c.id
// 		JOIN products p ON p.id = ci.product_id
// 		WHERE c.user_id = $1
