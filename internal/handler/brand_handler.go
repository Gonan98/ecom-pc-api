package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/service"
)

type BrandHandler struct {
	brandService *service.BrandService
}

func NewBrandHandler(brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

func (h *BrandHandler) Routes(r chi.Router) {
	r.Get("/", httpHandler(h.GetAll))
	r.Get("/{id}", httpHandler(h.GetByID))
}

func (h *BrandHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	brands, err := h.brandService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, brands)
}

func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return model.NewAPIError(http.StatusBadRequest, err)
	}

	brand, err := h.brandService.GetByID(r.Context(), ID)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, brand)
}
