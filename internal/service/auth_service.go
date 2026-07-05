package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/auth"
	"github.com/gonan98/ecom-pc-api/internal/errors"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/store"
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

	ok, err := s.userStore.ExistsUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if ok {
		return errors.NewAPIError(http.StatusBadRequest, fmt.Errorf("Email is already registered"))
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.userStore.Create(ctx, user, roleID)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !auth.ComparePasswords(user.Password, []byte(password)) {
		return "", err
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

func (s *AuthService) GetProfile(ctx context.Context) {

}
