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
	errInvalidJSON = types.NewAPIError(http.StatusBadRequest, errors.New("Invalid JSON payload"))
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
	r.Post("/register", httpHandler(h.register))
	r.Post("/login", httpHandler(h.login))

	r.With(middleware.JWTMiddleware).Get("/profile", httpHandler(h.getProfile))
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) error {
	var req types.CreateUserRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.service.Register(r.Context(), &req); err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusCreated, Message: "User created"})
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) error {
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

	return write(w, types.APIResponse{Code: http.StatusOK, Message: "Access token", Data: token})
}

func (h *AuthHandler) getProfile(w http.ResponseWriter, r *http.Request) error {
	userInfo, err := h.service.Profile(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: userInfo})
}
