package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/service"
)

type ShoppingCartHandler struct {
	shoppingCartService *service.ShoppingCartService
}

func NewShoppingCartHandler(shoppingCartService *service.ShoppingCartService) *ShoppingCartHandler {
	return &ShoppingCartHandler{
		shoppingCartService: shoppingCartService,
	}
}

func (h *ShoppingCartHandler) Routes(r chi.Router) {
	r.With(middleware.JWTMiddleware).Get("/", HttpHandler(h.Get))
}

func (h *ShoppingCartHandler) Get(w http.ResponseWriter, r *http.Request) error {
	cart, err := h.shoppingCartService.Get(r.Context())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, cart)
}
