package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/types"
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
	r.Use(middleware.JWTMiddleware)

	r.Post("/", httpHandler(h.create))
	r.Get("/", httpHandler(h.getOrders))
	r.Get("/{id}/details", httpHandler(h.getOrderDetails))
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) error {
	if err := h.orderService.Create(r.Context()); err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusCreated, Message: "Order created"})
}

func (h *OrderHandler) getOrders(w http.ResponseWriter, r *http.Request) error {
	orders, err := h.orderService.GetOrders(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: orders})
}

func (h *OrderHandler) getOrderDetails(w http.ResponseWriter, r *http.Request) error {
	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return types.NewAPIError(http.StatusBadRequest, err)
	}

	res, err := h.orderService.GetOrderItems(r.Context(), orderID)
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: res})
}
