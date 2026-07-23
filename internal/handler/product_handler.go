package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gonan98/ecom-pc-api/internal/middleware"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
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
	r.Get("/", httpHandler(h.getBrands))
	r.Get("/{id}", httpHandler(h.getBrandByID))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Post("/", httpHandler(h.createBrand))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Put("/{id}", httpHandler(h.updateBrand))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Delete("/{id}", httpHandler(h.deleteBrand))
}

func (h *ProductHandler) getBrands(w http.ResponseWriter, r *http.Request) error {
	products, err := h.productService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: products})
}

func (h *ProductHandler) getBrandByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	product, err := h.productService.GetByID(r.Context(), ID)
	if err != nil {
		return err
	}

	return writeResponse(w, types.APIResponse{Code: http.StatusOK, Data: product})
}

func (h *ProductHandler) createBrand(w http.ResponseWriter, r *http.Request) error {
	var req types.CreateProductRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.productService.Create(r.Context(), &req); err != nil {
		return err
	}

	return writeResponse(w, types.NewAPIResponse(http.StatusCreated, "New product created"))
}

func (h *ProductHandler) updateBrand(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	var req types.UpdateProductRequest

	if err := readJSON(r, &req); err != nil {
		return err
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.productService.Update(r.Context(), &req, ID); err != nil {
		return err
	}

	return writeResponse(w, types.NewAPIResponse(http.StatusCreated, "Product updated"))
}

func (h *ProductHandler) deleteBrand(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	if err := h.productService.Delete(r.Context(), ID); err != nil {
		return err
	}

	return writeResponse(w, types.NewAPIResponse(http.StatusCreated, "Product deleted"))
}
