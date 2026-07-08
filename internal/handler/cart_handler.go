package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/service"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) Routes(r chi.Router) {
	r.With(middleware.JWTMiddleware).Get("/", httpHandler(h.getCart))
	r.With(middleware.JWTMiddleware).Post("/items", httpHandler(h.addItem))
	r.With(middleware.JWTMiddleware).Delete("/items", httpHandler(h.deleteAllItems))
	r.With(middleware.JWTMiddleware).Delete("/items/{productID}", httpHandler(h.deleteItemByProductID))
	r.With(middleware.JWTMiddleware).Patch("/items/{productID}", httpHandler(h.updateItemQuantity))
}

func (h *CartHandler) getCart(w http.ResponseWriter, r *http.Request) error {
	cartResponse, err := h.cartService.GetCart(r.Context())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, cartResponse)
}

func (h *CartHandler) addItem(w http.ResponseWriter, r *http.Request) error {
	var req model.CartItemRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return model.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	err := h.cartService.AddItemToCart(r.Context(), &model.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{"message": "Item added to the cart"})
}

func (h *CartHandler) deleteAllItems(w http.ResponseWriter, r *http.Request) error {
	if err := h.cartService.DeleteCartItems(r.Context()); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *CartHandler) deleteItemByProductID(w http.ResponseWriter, r *http.Request) error {
	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		return model.NewAPIError(http.StatusBadRequest, fmt.Errorf("productID parameter must be an integer"))
	}

	if err := h.cartService.DeleteCartItemByProductID(r.Context(), productID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *CartHandler) updateItemQuantity(w http.ResponseWriter, r *http.Request) error {
	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		return model.NewAPIError(http.StatusBadRequest, fmt.Errorf("productID parameter must be an integer"))
	}

	var req model.UpdateCartItemRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return model.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	if err := h.cartService.UpdateItemQuantity(r.Context(), productID, req.Quantity); err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{"message": "Updated item quantity"})
}
