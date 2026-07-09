package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) Routes(r chi.Router) {
	r.With(middleware.JWTMiddleware).Post("/", httpHandler(h.create))
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) error {
	if err := h.orderService.Create(r.Context()); err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Order created",
	})
}
