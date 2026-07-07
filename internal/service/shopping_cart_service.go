package service

import (
	"context"
	"errors"

	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/repository"
)

type ShoppingCartService struct {
	shoppingCartRepo *repository.ShoppingCartRepository
}

func NewShoppingCartService(shoppingCartRepo *repository.ShoppingCartRepository) *ShoppingCartService {
	return &ShoppingCartService{
		shoppingCartRepo: shoppingCartRepo,
	}
}

func (s *ShoppingCartService) Get(ctx context.Context) (*model.ShoppingCart, error) {
	claims, err := middleware.GetUserClaims(ctx)
	if err != nil {
		return nil, errors.New("Error retrieving user claims from context")
	}

	userID, err := claims.UserID()
	if err != nil {
		return nil, errors.New("Error in parsing userID to int")
	}

	cart, err := s.shoppingCartRepo.Get(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	if cart.ID > 0 {
		return cart, nil
	}

	cartID, err := s.shoppingCartRepo.Create(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	cart.ID = cartID
	cart.UserID = int(userID)

	return cart, nil
}
