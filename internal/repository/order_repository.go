package repository

import (
	"context"

	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	// db *pgxpool.Pool
	db DBTX
}

func NewOrderRepository(db DBTX) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) WithTx(tx pgx.Tx) *OrderRepository {
	return &OrderRepository{
		db: tx,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *types.Order) (int, error) {
	orderID := 0
	query := "INSERT INTO orders (user_id, total) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(ctx, query, order.UserID, order.Total).Scan(&orderID)
	return orderID, err
}

func (r *OrderRepository) CreateDetail(ctx context.Context, orderDetail *types.OrderDetail) error {
	query := "INSERT INTO order_details (order_id, product_id, quantity, unit_price, discount) VALUES ($1,$2,$3,$4,$5)"
	_, err := r.db.Exec(ctx, query, orderDetail.OrderID, orderDetail.ProductID, orderDetail.Quantity, orderDetail.UnitPrice, orderDetail.Discount)
	return err
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]types.Order, error) {
	query := "SELECT id, user_id, status, total, created_at FROM orders"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0)
	for rows.Next() {
		var o types.Order
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

func (r *OrderRepository) GetByUser(ctx context.Context, userID int) ([]types.Order, error) {
	query := "SELECT id, user_id, status, total, created_at FROM orders WHERE user_id = $1"
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0)
	for rows.Next() {
		var o types.Order
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

func (r *OrderRepository) GetByID(ctx context.Context, orderID int) (*types.Order, error) {
	order := new(types.Order)
	query := "SELECT id, user_id, status, total, created_at FROM orders WHERE id = $1"
	err := r.db.QueryRow(ctx, query, orderID).Scan(&order.ID, &order.UserID, &order.Status, &order.Total, &order.CreatedAt)

	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	return order, nil
}

// func (r *OrderRepository) ExistOrderByIDAndUserID(ctx context.Context, orderID int, userID int) (bool, error) {
// 	result := false
// 	query := "SELECT EXISTS (SELECT 1 FROM orders WHERE id = $1 AND user_id = $2)"
// 	err := r.db.QueryRow(ctx, query, orderID, userID).Scan(&result)
// 	return result, err
// }

func (r *OrderRepository) GetDetailsByOrder(ctx context.Context, orderID int) ([]types.OrderDetail, error) {
	query := "SELECT od.order_id, od.product_id, od.quantity, od.unit_price, od.discount FROM order_details od JOIN orders o ON o.id = od.order_id WHERE od.order_id = $1"
	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	details := make([]types.OrderDetail, 0)
	for rows.Next() {
		var od types.OrderDetail
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

func (r *OrderRepository) GetDetailsByOrderAndUser(ctx context.Context, userID int, orderID int) ([]types.OrderDetail, error) {
	query := "SELECT od.order_id, od.product_id, od.quantity, od.unit_price, od.discount FROM order_details od JOIN orders o ON o.id = od.order_id WHERE o.user_id = $1 AND od.order_id = $2"
	rows, err := r.db.Query(ctx, query, userID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	details := make([]types.OrderDetail, 0)
	for rows.Next() {
		var od types.OrderDetail
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

func (r *OrderRepository) UpdateStatus(ctx context.Context, status string, orderID int) error {
	query := "UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
	_, err := r.db.Exec(ctx, query, status, orderID)
	return err
}
