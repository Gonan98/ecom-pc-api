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

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) Routes(r chi.Router) {
	r.Use(middleware.JWTMiddleware)

	r.Get("/", httpHandler(h.getCart))
	r.Post("/items", httpHandler(h.addItem))
	r.Delete("/items", httpHandler(h.deleteAllItems))
	r.Delete("/items/{productID}", httpHandler(h.deleteItemByProductID))
	r.Patch("/items/{productID}", httpHandler(h.updateItemQuantity))
}

func (h *CartHandler) getCart(w http.ResponseWriter, r *http.Request) error {
	cartResponse, err := h.cartService.GetCart(r.Context())
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: cartResponse})
}

func (h *CartHandler) addItem(w http.ResponseWriter, r *http.Request) error {
	var req types.CartItemRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	err := h.cartService.AddItem(r.Context(), &types.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		return err
	}

	return writeResponse(w, types.NewAPIResponse(http.StatusOK, fmt.Sprintf("Product %d added to cart", req.ProductID)))
}

func (h *CartHandler) deleteAllItems(w http.ResponseWriter, r *http.Request) error {
	if err := h.cartService.DeleteItems(r.Context()); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *CartHandler) deleteItemByProductID(w http.ResponseWriter, r *http.Request) error {
	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		return util.InvalidParamID("productID")
	}

	if err := h.cartService.DeleteItemByProductID(r.Context(), productID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *CartHandler) updateItemQuantity(w http.ResponseWriter, r *http.Request) error {
	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		return util.InvalidParamID("productID")
	}

	var req types.UpdateCartItemRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.cartService.UpdateItemQuantity(r.Context(), productID, req.Quantity); err != nil {
		return err
	}

	return writeResponse(w, types.NewAPIResponse(http.StatusOK, fmt.Sprintf("Item quantity updated for product %d", productID)))
}
