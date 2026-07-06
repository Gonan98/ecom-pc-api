package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: service,
	}
}

func (h *ProductHandler) Routes(r chi.Router) {
	r.Get("/", HttpHandler(h.GetAll))
	r.Get("/{id}", HttpHandler(h.GetByID))
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	products, err := h.productService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return model.NewAPIError(http.StatusBadRequest, err)
	}

	product, err := h.productService.GetByID(r.Context(), ID)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, product)
}
