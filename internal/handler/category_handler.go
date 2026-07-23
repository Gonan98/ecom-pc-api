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

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Routes(r chi.Router) {
	r.Get("/", httpHandler(h.getCategories))
	r.Get("/{id}", httpHandler(h.getCategoryByID))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Post("/", httpHandler(h.createCategory))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Put("/{id}", httpHandler(h.updateCategory))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Delete("/{id}", httpHandler(h.deleteCategory))
}

func (h *CategoryHandler) getCategories(w http.ResponseWriter, r *http.Request) error {
	cateogries, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: cateogries})
}

func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return types.NewAPIError(http.StatusBadRequest, err)
	}

	category, err := h.categoryService.GetByID(r.Context(), ID)
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: category})
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateCategoryRequest)

	if err := readJSON(r, req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.categoryService.Create(r.Context(), req); err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusCreated, Message: "Category created"})
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	req := new(types.UpdateCategoryRequest)

	if err := readJSON(r, req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.categoryService.Update(r.Context(), req, ID); err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Message: fmt.Sprintf("Category with ID = %d updated", ID)})
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	if err := h.categoryService.Delete(r.Context(), ID); err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Message: fmt.Sprintf("Category with ID = %d deleted", ID)})
}
