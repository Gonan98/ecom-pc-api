package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
)

var (
	errInvalidJSON = types.NewAPIError(http.StatusBadRequest, errors.New("Invalid JSON structure"))
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
	r.Post("/register", httpHandler(h.Register))
	r.Post("/login", httpHandler(h.Login))

	r.With(middleware.JWTMiddleware).Get("/profile", httpHandler(h.Profile))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req types.CreateUserRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	err := h.service.Register(r.Context(), types.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password,
	})

	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusCreated, Message: "User created"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req types.LogUserRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	token, err := h.service.Login(r.Context(), &req)
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Message: "Receiving token", Data: token})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) error {
	u, err := h.service.Profile(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: u})
}
