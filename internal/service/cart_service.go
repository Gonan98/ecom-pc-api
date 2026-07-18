package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

var (
	errCartIsEmpty           = types.NewAPIError(http.StatusBadRequest, errors.New("cart is empty"))
	errProductNotFoundInCart = types.NewAPIError(http.StatusBadRequest, errors.New("product is not in the cart"))
	errProductAlreadyInCart  = types.NewAPIError(http.StatusBadRequest, errors.New("that product is already in the cart"))
)

type CartService struct {
	cartRepo    *repo.CartRepository
	productRepo *repo.ProductRepository
}

func NewCartService(cartRepo *repo.CartRepository, productRepo *repo.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) AddItem(ctx context.Context, cartItem *types.CartItem) error {

	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	ok, err := s.cartRepo.ExistsItemInCartByProductID(ctx, cart.ID, cartItem.ProductID)
	if err != nil {
		return err
	}

	if ok {
		return errProductAlreadyInCart
	}

	cartItem.CartID = cart.ID

	return s.cartRepo.CreateItem(ctx, cartItem)
}

func (s *CartService) GetCart(ctx context.Context) (*types.CartResponse, error) {

	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	cartItems, err := s.cartRepo.GetItemsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.cartToResponse(ctx, cartItems)
}

func (s *CartService) DeleteItems(ctx context.Context) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetByUser(ctx, userID)
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

	return s.cartRepo.DeleteItems(ctx, cart.ID)
}

func (s *CartService) DeleteItemByProductID(ctx context.Context, productID int) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetByUser(ctx, userID)
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

	return s.cartRepo.DeleteItemsByProductID(ctx, cart.ID, productID)
}

func (s *CartService) UpdateItemQuantity(ctx context.Context, productID int, quantity int) error {
	userID, _, err := extractUserFromClaims(ctx)
	if err != nil {
		return err
	}

	cart, err := s.cartRepo.GetByUser(ctx, userID)
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

func (s *CartService) cartToResponse(ctx context.Context, cartItems []types.CartItem) (*types.CartResponse, error) {
	cartResponse := new(types.CartResponse)
	cartResponse.Items = make([]types.CartItemResponse, 0)
	for _, item := range cartItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, types.NewAPIError(http.StatusBadRequest, fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID))
		}

		itemResp := types.CartItemResponse{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			Subtotal:    float64(item.Quantity) * product.Price,
		}

		cartResponse.Total += itemResp.Subtotal
		cartResponse.Items = append(cartResponse.Items, itemResp)
	}

	return cartResponse, nil
}
