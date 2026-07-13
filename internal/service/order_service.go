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
		productIDs := make(map[int]*types.Product)

		for _, item := range cartItems {
			product, err := s.productRepo.GetByID(ctx, item.ProductID)
			if err != nil {
				return err
			}

			if item.Quantity > product.Stock {
				return types.NewAPIError(http.StatusBadRequest, fmt.Errorf("Product %s is not available in the quantity requested", product.Name))
			}

			productIDs[item.ProductID] = product
			total += product.Price * float64(item.Quantity)
		}

		orderID, err := orderTx.Create(ctx, &types.Order{
			UserID: userID,
			Total:  total,
		})

		if err != nil {
			return err
		}

		for _, item := range cartItems {
			product := productIDs[item.ProductID]

			if err := orderTx.CreateDetail(ctx, &types.OrderDetail{
				OrderID:   orderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: product.Price,
				Discount:  0,
			}); err != nil {
				return err
			}

			product.Stock -= item.Quantity

			if err := productTx.UpdateStock(ctx, item.ProductID, product.Stock); err != nil {
				return err
			}
		}

		if err := cartTx.DeleteCartItems(ctx, cart.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderService) GetOrders(ctx context.Context) ([]types.Order, error) {

	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	return s.orderRepo.GetOrders(ctx, userID)
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID int) ([]types.OrderDetailResponse, error) {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	details, err := s.orderRepo.GetOrderDetails(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	if len(details) == 0 {
		return nil, types.NewAPIError(http.StatusNotFound, fmt.Errorf("you don't have an order with ID=%d", orderID))
	}

	detailsResponse := make([]types.OrderDetailResponse, 0)
	for _, detail := range details {
		p, err := s.productRepo.GetByID(ctx, detail.ProductID)
		if err != nil {
			return nil, err
		}

		dr := types.OrderDetailResponse{
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
