package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/database"
	"github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/jackc/pgx/v5"
)

var (
	errInvalidEmailOrPassword = types.NewAPIError(http.StatusBadRequest, errors.New("Invalid email or password"))
	errEmailAlreadyRegistered = types.NewAPIError(http.StatusBadRequest, errors.New("Email is already registered"))
)

type AuthService struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	cartRepo  *repository.CartRepository
	txManager *database.TxManager
}

func NewAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	cartRepo *repository.CartRepository,
	txManager *database.TxManager,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		cartRepo:  cartRepo,
		txManager: txManager,
	}
}

func (s *AuthService) Register(ctx context.Context, user types.User) error {
	return s.txManager.RunInTx(ctx, func(tx pgx.Tx) error {

		userTx := s.userRepo.WithTx(tx)
		cartTx := s.cartRepo.WithTx(tx)

		roleID, err := s.roleRepo.GetCustomerRoleID(ctx)
		if err != nil {
			return fmt.Errorf("Role customer does not exist")
		}

		ok, err := s.userRepo.ExistByEmail(ctx, user.Email)
		if err != nil {
			return err
		}

		if ok {
			return errEmailAlreadyRegistered
		}

		hashedPassword, err := auth.HashPassword(user.PasswordHash)
		if err != nil {
			return err
		}
		user.PasswordHash = hashedPassword

		userID, err := userTx.Create(ctx, user, roleID)
		if err != nil {
			return err
		}

		if err := cartTx.Create(ctx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errInvalidEmailOrPassword
	}

	if !auth.ComparePasswords(user.PasswordHash, []byte(password)) {
		return "", errInvalidEmailOrPassword
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateJWT(user.ID, role.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Profile(ctx context.Context) (*types.UserInfo, error) {
	userID, role, err := extractUserFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo := &types.UserInfo{
		ID:        user.ID,
		RoleName:  role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return userInfo, nil
}

func extractUserFromClaims(ctx context.Context) (int, string, error) {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return 0, "", err
	}

	userID, err := claims.UserID()
	if err != nil {
		return 0, "", err
	}

	return userID, claims.Role, err
}
