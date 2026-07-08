package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/repository"
)

var (
	errCartIsEmpty           = model.NewAPIError(http.StatusBadRequest, errors.New("cart is empty"))
	errProductNotFoundInCart = model.NewAPIError(http.StatusBadRequest, errors.New("product is not in the cart"))
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

func (s *CartService) AddItemToCart(ctx context.Context, cartItem *model.CartItem) error {

	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	cartItem.CartID = cart.ID

	return s.cartRepo.CreateItem(ctx, cartItem)
}

func (s *CartService) GetCart(ctx context.Context) (*model.CartResponse, error) {

	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	cartItems, err := s.cartRepo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	var cartResponse model.CartResponse

	for _, item := range cartItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		itemResp := model.CartItemResponse{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			Subtotal:    float64(item.Quantity) * product.Price,
		}

		cartResponse.Total += itemResp.Subtotal
		cartResponse.Items = append(cartResponse.Items, itemResp)
	}

	return &cartResponse, nil
}

func (s *CartService) DeleteCartItems(ctx context.Context) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	ok, err := s.cartRepo.ExistsItemsInCart(ctx, cart.ID)
	if err != nil {
		return err
	}

	if !ok {
		return errCartIsEmpty
	}

	return s.cartRepo.DeleteCartItems(ctx, cart.ID)
}

func (s *CartService) DeleteCartItemByProductID(ctx context.Context, productID int) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	ok, err := s.cartRepo.ExistsItemInCartByProductID(ctx, cart.ID, productID)
	if err != nil {
		return err
	}

	if !ok {
		return errProductNotFoundInCart
	}

	return s.cartRepo.DeleteCartItemsByProductID(ctx, cart.ID, productID)
}

func (s *CartService) UpdateItemQuantity(ctx context.Context, productID int, quantity int) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	ok, err := s.cartRepo.ExistsItemInCartByProductID(ctx, cart.ID, productID)
	if err != nil {
		return err
	}

	if !ok {
		return errProductNotFoundInCart
	}

	return s.cartRepo.UpdateItemQuantity(ctx, cart.ID, productID, quantity)
}
