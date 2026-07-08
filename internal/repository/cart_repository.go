package repository

import (
	"context"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) Create(ctx context.Context, userID int) error {
	q := "INSERT INTO shopping_carts (user_id) VALUES ($1)"
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

func (r *CartRepository) CreateItem(ctx context.Context, item *model.CartItem) error {
	q := "INSERT INTO shopping_cart_items VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, q, item.CartID, item.ProductID, item.Quantity)
	return err
}

func (r *CartRepository) GetCart(ctx context.Context, userID int) (*model.Cart, error) {
	var cart model.Cart
	q := "SELECT id, user_id FROM shopping_carts WHERE user_id = $1"
	err := r.db.QueryRow(ctx, q, userID).Scan(&cart.ID, &cart.UserID)

	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	return &cart, nil
}

func (r *CartRepository) GetCartItems(ctx context.Context, userID int) ([]model.CartItem, error) {
	q := "SELECT ci.product_id, ci.quantity FROM shopping_cart_items ci JOIN shopping_carts c ON ci.cart_id = c.id WHERE c.user_id = $1"
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.CartItem, 0)
	for rows.Next() {
		var item model.CartItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
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

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *CartRepository) ExistsItemInCartByProductID(ctx context.Context, cartID int, productID int) (bool, error) {
	result := false
	q := "SELECT EXISTS (SELECT 1 FROM shopping_cart_items WHERE cart_id = $1 AND product_id = $2)"
	err := r.db.QueryRow(ctx, q, cartID, productID).Scan(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *CartRepository) DeleteCartItems(ctx context.Context, cartID int) error {
	q := "DELETE FROM shopping_cart_items WHERE cart_id = $1"
	_, err := r.db.Exec(ctx, q, cartID)
	return err
}

func (r *CartRepository) DeleteCartItemsByProductID(ctx context.Context, cartID int, productID int) error {
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
