package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/repository"
)

var (
	errInvalidEmailOrPassword = model.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid email or password"))
)

type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, user model.User) error {
	roleID, err := s.roleRepo.GetCustomerRoleID(ctx)
	if err != nil {
		return fmt.Errorf("Role customer does not exist")
	}

	ok, err := s.userRepo.ExistByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if ok {
		return model.NewAPIError(http.StatusBadRequest, fmt.Errorf("Email is already registered"))
	}

	hashedPassword, err := auth.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	return s.userRepo.Create(ctx, user, roleID)
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

func (s *AuthService) Profile(ctx context.Context) (*model.User, error) {
	claims, err := middleware.GetUserClaims(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := claims.UserID()
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
