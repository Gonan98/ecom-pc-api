package handler

import (
	"net/http"

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
	r.With(middleware.JWTMiddleware).Get("/", HttpHandler(h.getCart))
	r.With(middleware.JWTMiddleware).Post("/items", HttpHandler(h.addItem))
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

	err := h.cartService.CreateItem(r.Context(), &model.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{"message": "Item added to the cart"})
}
