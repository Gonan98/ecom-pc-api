package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/database"
	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
	"github.com/jackc/pgx/v5"
)

type OrderService struct {
	orderRepo   *repo.OrderRepository
	productRepo *repo.ProductRepository
	cartRepo    *repo.CartRepository
	txManager   *database.TxManager
}

func NewOrderService(
	orderRepo *repo.OrderRepository,
	productRepo *repo.ProductRepository,
	cartRepo *repo.CartRepository,
	txManager *database.TxManager,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		txManager:   txManager,
	}
}

func (s *OrderService) Create(ctx context.Context) error {
	return s.txManager.RunInTx(ctx, func(tx pgx.Tx) error {

		orderTx := s.orderRepo.WithTx(tx)
		productTx := s.productRepo.WithTx(tx)
		cartTx := s.cartRepo.WithTx(tx)

		userID, _, err := extractUserFromClaims(ctx)
		if err != nil {
			return err
		}

		cart, err := cartTx.GetByUser(ctx, userID)
		if err != nil {
			return err
		}

		cartItems, err := cartTx.GetItemsByUser(ctx, userID)
		if err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return errCartIsEmpty
		}

		var total float64
		prices := make(map[int]float64)

		// Calculate total from cart and decrease stock
		for _, item := range cartItems {
			product, err := productTx.GetByID(ctx, item.ProductID)
			if err != nil {
				return err
			}

			if item.Quantity > product.Stock {
				return util.NotAvailableStock(product.Name)
			}

			prices[item.ProductID] = product.Price
			total += product.Price * float64(item.Quantity) * (1 - item.Discount)
			if err := productTx.DecreaseStock(ctx, item.Quantity, product.ID); err != nil {
				return err
			}
		}

		// Create an Order
		orderID, err := orderTx.Create(ctx, &types.Order{
			UserID: userID,
			Total:  total,
		})

		if err != nil {
			return err
		}

		// Create OrderDetails
		for _, item := range cartItems {
			price := prices[item.ProductID]

			err := orderTx.CreateDetail(ctx, &types.OrderDetail{
				OrderID:   orderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: price,
				Discount:  item.Discount,
			})

			if err != nil {
				return err
			}
		}

		// Clean cart
		if err := cartTx.DeleteItems(ctx, cart.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderService) GetOrders(ctx context.Context) ([]types.Order, error) {
	userID, role, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	if role == types.RoleNameAdmin {
		return s.orderRepo.GetAll(ctx)
	}

	return s.orderRepo.GetByUser(ctx, userID)
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID int) ([]types.OrderDetailResponse, error) {
	userID, role, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	if role == types.RoleNameAdmin {
		details, err := s.orderRepo.GetDetailsByOrder(ctx, orderID)
		if err != nil {
			return nil, err
		}

		if len(details) == 0 {
			return nil, types.NewAPIError(http.StatusNotFound, fmt.Errorf("order with ID: %d not found", orderID))
		}

		return s.orderDetailToResponse(ctx, details)
	}

	details, err := s.orderRepo.GetDetailsByOrderAndUser(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	if len(details) == 0 {
		return nil, types.NewAPIError(http.StatusNotFound, fmt.Errorf("you don't have an order with ID: %d", orderID))
	}

	return s.orderDetailToResponse(ctx, details)
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID int, status types.OrderStatus) error {
	return s.txManager.RunInTx(ctx, func(tx pgx.Tx) error {
		orderTx := s.orderRepo.WithTx(tx)
		productTx := s.productRepo.WithTx(tx)

		order, err := orderTx.GetByID(ctx, orderID)
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewAPIError(http.StatusNotFound, fmt.Errorf("order with ID=%d does not exist", orderID))
		}

		if err != nil {
			return err
		}

		if !s.isValidTransition(order.Status, status) {
			return types.NewAPIError(http.StatusBadRequest, fmt.Errorf("cannot transition order status from %s to %s", order.Status, status))
		}

		if status == types.OrderStatusCancelled {
			details, err := orderTx.GetDetailsByOrder(ctx, orderID)
			if err != nil {
				return err
			}

			for _, d := range details {
				if err := productTx.IncreaseStock(ctx, d.Quantity, d.ProductID); err != nil {
					return err
				}
			}
		}

		return orderTx.UpdateStatus(ctx, status, orderID)
	})
}

func (s *OrderService) orderDetailToResponse(ctx context.Context, details []types.OrderDetail) ([]types.OrderDetailResponse, error) {
	response := make([]types.OrderDetailResponse, 0)
	for _, od := range details {
		p, err := s.productRepo.GetByID(ctx, od.ProductID)
		if err != nil {
			return nil, err
		}

		dr := types.OrderDetailResponse{
			ProductID:   od.ProductID,
			ProductName: p.Name,
			UnitPrice:   od.UnitPrice,
			Quantity:    od.Quantity,
			Discount:    od.Discount,
		}

		response = append(response, dr)
	}

	return response, nil
}

func (s *OrderService) isValidTransition(currentStatus types.OrderStatus, newStatus types.OrderStatus) bool {
	return (currentStatus == types.OrderStatusPending && (newStatus == types.OrderStatusPaid || newStatus == types.OrderStatusCancelled)) ||
		(currentStatus == types.OrderStatusPaid && (newStatus == types.OrderStatusShipped || newStatus == types.OrderStatusCancelled)) ||
		(currentStatus == types.OrderStatusShipped && newStatus == types.OrderStatusDelivered)
}
