package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) CreateItem(ctx context.Context, cartItem *model.CartItem) error {

	userID, err := extractUserIDFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart.ID == 0 {
		cart.ID, err = s.cartRepo.Create(ctx, userID)
		if err != nil {
			return err
		}
	}

	cartItem.CartID = cart.ID

	return s.cartRepo.CreateItem(ctx, cartItem)
}

func (s *CartService) GetCart(ctx context.Context) (*model.CartResponse, error) {

	userID, err := extractUserIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart.ID == 0 {
		return nil, model.NewAPIError(http.StatusNotFound, fmt.Errorf("You don't have a cart, start adding a product"))
	}

	cartItems, err := s.cartRepo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	var cartResponse model.CartResponse
	cartResponse.Items = make([]model.CartItemResponse, 0)

	for _, item := range cartItems {
		var itemResponse model.CartItemResponse
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		itemResponse.ProductID = product.ID
		itemResponse.ProductName = product.Name
		itemResponse.Quantity = item.Quantity
		itemResponse.UnitPrice = product.Price
		itemResponse.Subtotal = float64(item.Quantity) * product.Price

		cartResponse.Total += itemResponse.Subtotal
		cartResponse.Items = append(cartResponse.Items, itemResponse)
	}

	return &cartResponse, nil
}
