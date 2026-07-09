package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	cartRepo    *repository.CartRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	cartRepo *repository.CartRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
	}
}

func (s *OrderService) Create(ctx context.Context) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	cartItems, err := s.cartRepo.GetCartItems(ctx, userID)
	if err != nil {
		return err
	}

	if len(cartItems) == 0 {
		return errCartIsEmpty
	}

	var total float64
	productIDs := make(map[int]*model.Product)

	for _, item := range cartItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return err
		}

		if item.Quantity > product.Stock {
			return model.NewAPIError(http.StatusBadRequest, fmt.Errorf("Product %s is not available in the quantity requested", product.Name))
		}

		productIDs[item.ProductID] = product
		total += product.Price * float64(item.Quantity)
	}

	orderID, err := s.orderRepo.Create(ctx, &model.Order{
		UserID: userID,
		Total:  total,
	})

	if err != nil {
		return err
	}

	for _, item := range cartItems {
		product := productIDs[item.ProductID]

		if err := s.orderRepo.CreateDetail(ctx, &model.OrderDetail{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
			Discount:  0,
		}); err != nil {
			return err
		}

		product.Stock -= item.Quantity

		if err := s.productRepo.UpdateStock(ctx, item.ProductID, product.Stock); err != nil {
			return err
		}
	}

	return s.cartRepo.DeleteCartItems(ctx, cart.ID)
}

func (s *OrderService) GetOrders(ctx context.Context) ([]model.Order, error) {

	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	return s.orderRepo.GetOrders(ctx, userID)
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID int) ([]model.OrderDetailResponse, error) {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	details, err := s.orderRepo.GetOrderDetails(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	if len(details) == 0 {
		return nil, model.NewAPIError(http.StatusNotFound, fmt.Errorf("you don't have an order with ID=%d", orderID))
	}

	detailsResponse := make([]model.OrderDetailResponse, 0)
	for _, detail := range details {
		p, err := s.productRepo.GetByID(ctx, detail.ProductID)
		if err != nil {
			return nil, err
		}

		dr := model.OrderDetailResponse{
			ProductID:   detail.ProductID,
			ProductName: p.Name,
			UnitPrice:   detail.UnitPrice,
			Quantity:    detail.Quantity,
			Discount:    detail.Discount,
		}

		detailsResponse = append(detailsResponse, dr)
	}

	return detailsResponse, nil
}
