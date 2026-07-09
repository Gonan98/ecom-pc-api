package repository

import (
	"context"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) (int, error) {
	orderID := 0
	query := "INSERT INTO orders (user_id, total) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(ctx, query, order.UserID, order.Total).Scan(&orderID)
	return orderID, err
}

func (r *OrderRepository) CreateDetail(ctx context.Context, orderDetail *model.OrderDetail) error {
	query := "INSERT INTO order_details (order_id, product_id, quantity, unit_price, discount) VALUES ($1,$2,$3,$4,$5)"
	_, err := r.db.Exec(ctx, query, orderDetail.OrderID, orderDetail.ProductID, orderDetail.Quantity, orderDetail.UnitPrice, orderDetail.Discount)
	return err
}
