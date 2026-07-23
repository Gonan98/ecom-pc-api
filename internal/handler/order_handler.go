package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
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

	r.Post("/", httpHandler(h.createOrder))
	r.Get("/", httpHandler(h.getOrders))
	r.Get("/{id}/details", httpHandler(h.getOrderDetails))

	// r.With(middleware.AdminMiddleware).Get("/", httpHandler(h.getAll))
	r.With(middleware.AdminMiddleware).Patch("/{id}/status", httpHandler(h.updateStatus))
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) error {
	if err := h.orderService.Create(r.Context()); err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusCreated, Message: "Order created"})
}

func (h *OrderHandler) getOrders(w http.ResponseWriter, r *http.Request) error {
	orders, err := h.orderService.GetOrders(r.Context())
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: orders})
}

func (h *OrderHandler) getOrderDetails(w http.ResponseWriter, r *http.Request) error {
	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	res, err := h.orderService.GetOrderItems(r.Context(), orderID)
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: res})
}

func (h *OrderHandler) updateStatus(w http.ResponseWriter, r *http.Request) error {
	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	var req types.UpdateOrderStatusRequest
	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.orderService.UpdateStatus(r.Context(), orderID, req.Status); err != nil {
		return err
	}

	return writeResponse(w, types.NewAPIResponse(http.StatusOK, fmt.Sprintf("Order %d status updated", orderID)))
}
