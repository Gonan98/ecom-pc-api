package service

import (
	"context"
	"errors"
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
	cartRepo       *repo.CartRepository
	productService *ProductService
}

func NewCartService(cartRepo *repo.CartRepository, productService *ProductService) *CartService {
	return &CartService{
		cartRepo:       cartRepo,
		productService: productService,
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

	_, err = s.productService.GetByID(ctx, cartItem.ProductID)
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

	cartItems, err := s.cartRepo.GetItemsWithProductsByUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	return s.cartToResponse(cartItems), nil
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

func (s *CartService) cartToResponse(items []repo.CartItemWithProduct) *types.CartResponse {
	resp := types.CartResponse{
		Total: 0,
		Items: make([]types.CartItemResponse, 0),
	}

	for _, item := range items {
		itemResp := types.CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Discount:    item.Discount,
			Subtotal:    float64(item.Quantity) * item.UnitPrice * (1 - item.Discount),
		}

		resp.Total += itemResp.Subtotal
		resp.Items = append(resp.Items, itemResp)
	}

	return &resp
}
