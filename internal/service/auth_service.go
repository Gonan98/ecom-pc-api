package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/errors"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/store"
)

var (
	errInvalidEmailOrPassword = errors.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid email or password"))
)

type AuthService struct {
	userStore *store.UserStore
	roleStore *store.RoleStore
}

func NewAuthService(userStore *store.UserStore, roleStore *store.RoleStore) *AuthService {
	return &AuthService{
		userStore: userStore,
		roleStore: roleStore,
	}
}

func (s *AuthService) Register(ctx context.Context, user model.User) error {
	roleID, err := s.roleStore.GetCustomerRoleID(ctx)
	if err != nil {
		return fmt.Errorf("Role customer does not exist")
	}

	ok, err := s.userStore.ExistByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if ok {
		return errors.NewAPIError(http.StatusBadRequest, fmt.Errorf("Email is already registered"))
	}

	hashedPassword, err := auth.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	return s.userStore.Create(ctx, user, roleID)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return "", errInvalidEmailOrPassword
	}

	log.Printf("User Password: %s \nRequest Password: %s", user.PasswordHash, password)

	if !auth.ComparePasswords(user.PasswordHash, []byte(password)) {
		return "", errInvalidEmailOrPassword
	}

	role, err := s.roleStore.GetByID(ctx, user.RoleID)
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

	user, err := s.userStore.GetByID(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
