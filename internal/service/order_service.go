package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/database"
	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
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

		cart, err := s.cartRepo.GetByUser(ctx, userID)
		if err != nil {
			return err
		}

		cartItems, err := s.cartRepo.GetItemsByUser(ctx, userID)
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
			product, err := s.productRepo.GetByID(ctx, item.ProductID)
			if err != nil {
				return err
			}

			if item.Quantity > product.Stock {
				return types.NewAPIError(http.StatusBadRequest, fmt.Errorf("Product %s is not available in the quantity requested", product.Name))
			}

			prices[item.ProductID] = product.Price
			total += product.Price * float64(item.Quantity)
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
				Discount:  0,
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

	if role == "admin" {
		return s.orderRepo.GetAll(ctx)
	}

	return s.orderRepo.GetByUser(ctx, userID)
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID int) ([]types.OrderDetailResponse, error) {
	userID, role, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	if role == "admin" {
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
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.ID == 0 {
		return types.NewAPIError(http.StatusNotFound, fmt.Errorf("order with ID=%d does not exist", orderID))
	}
	valid := (order.Status == string(types.OrderStatusPending) &&
		(status == types.OrderStatusPaid || status == types.OrderStatusCancelled)) ||
		(order.Status == string(types.OrderStatusPaid) &&
			(status == types.OrderStatusShipped || status == types.OrderStatusCancelled)) ||
		(order.Status == string(types.OrderStatusShipped) && status == types.OrderStatusDelivered)
	if !valid {
		return types.NewAPIError(http.StatusBadRequest, fmt.Errorf("cannot transition order status from %q to %q", order.Status, status))
	}

	return s.orderRepo.UpdateStatus(ctx, string(status), orderID)
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
