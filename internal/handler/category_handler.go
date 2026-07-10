package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Routes(r chi.Router) {
	r.Get("/", httpHandler(h.GetAll))
	r.Get("/{id}", httpHandler(h.GetByID))
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	cateogries, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, cateogries)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return types.NewAPIError(http.StatusBadRequest, err)
	}

	category, err := h.categoryService.GetByID(r.Context(), ID)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, category)
}
