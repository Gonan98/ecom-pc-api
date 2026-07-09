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

func (r *OrderRepository) GetOrders(ctx context.Context, userID int) ([]model.Order, error) {
	query := "SELECT id, user_id, status, total, created_at FROM orders WHERE user_id = $1"
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderDetails(ctx context.Context, userID int, orderID int) ([]model.OrderDetail, error) {
	query := "SELECT od.order_id, od.product_id, od.quantity, od.unit_price, od.discount FROM order_details od JOIN orders o ON o.id = od.order_id WHERE o.user_id = $1 AND od.order_id = $2"
	rows, err := r.db.Query(ctx, query, userID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	details := make([]model.OrderDetail, 0)
	for rows.Next() {
		var od model.OrderDetail
		if err := rows.Scan(&od.OrderID, &od.ProductID, &od.Quantity, &od.UnitPrice, &od.Discount); err != nil {
			return nil, err
		}
		details = append(details, od)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return details, nil
}
