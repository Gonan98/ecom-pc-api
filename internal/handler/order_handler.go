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

	r.Post("/", httpHandler(h.create))
	r.Get("/", httpHandler(h.getByUser))
	r.Get("/{id}/details", httpHandler(h.getDetails))

	r.With(middleware.AdminMiddleware).Get("/", httpHandler(h.getAll))
	r.With(middleware.AdminMiddleware).Patch("/{id}/status", httpHandler(h.updateStatus))
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) error {
	if err := h.orderService.Create(r.Context()); err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusCreated, Message: "Order created"})
}

func (h *OrderHandler) getAll(w http.ResponseWriter, r *http.Request) error {
	orders, err := h.orderService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: orders})
}

func (h *OrderHandler) getByUser(w http.ResponseWriter, r *http.Request) error {
	orders, err := h.orderService.GetAllByUser(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: orders})
}

func (h *OrderHandler) getDetails(w http.ResponseWriter, r *http.Request) error {
	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	res, err := h.orderService.GetOrderItems(r.Context(), orderID)
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: res})
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

	return write(w, types.NewAPIResponse(http.StatusOK, fmt.Sprintf("Order %d status updated", orderID)))
}
