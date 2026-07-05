package handler

import (
	// "log"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/errors"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/util"
)

var (
	errInvalidJSON = errors.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid json structure"))
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Post("/register", HttpHandler(h.Register))
	r.Post("/login", HttpHandler(h.Login))

	r.With(middleware.JWTMiddleware).Get("/profile", HttpHandler(h.Profile))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req model.CreateUserRequest

	if err := util.ReadJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return errors.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	err := h.service.Register(r.Context(), model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password,
	})

	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req model.LogUserRequest

	if err := util.ReadJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return errors.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"accessToken": token})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) error {
	u, err := h.service.Profile(r.Context())
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, u)
}
