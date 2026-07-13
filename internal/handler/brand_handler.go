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

type BrandHandler struct {
	brandService *service.BrandService
}

func NewBrandHandler(brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

func (h *BrandHandler) Routes(r chi.Router) {
	r.Get("/", httpHandler(h.getAll))
	r.Get("/{id}", httpHandler(h.getByID))

	r.With(middleware.JWTMiddleware, middleware.AdminMiddleware).Post("/", httpHandler(h.create))
	r.With(middleware.JWTMiddleware, middleware.AdminMiddleware).Put("/{id}", httpHandler(h.update))
	r.With(middleware.JWTMiddleware, middleware.AdminMiddleware).Delete("/{id}", httpHandler(h.delete))
}

func (h *BrandHandler) getAll(w http.ResponseWriter, r *http.Request) error {
	brands, err := h.brandService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: brands})
}

func (h *BrandHandler) getByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	brand, err := h.brandService.GetByID(r.Context(), ID)
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: brand})
}

func (h *BrandHandler) create(w http.ResponseWriter, r *http.Request) error {
	var req types.CreateBrandRequest

	if err := readJSON(r, &req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.brandService.Create(r.Context(), &req); err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Message: "New brand created"})
}

func (h *BrandHandler) update(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	var req *types.UpdateBrandRequest

	if err := readJSON(r, req); err != nil {
		return errInvalidJSON
	}

	if err := validate.Struct(req); err != nil {
		return util.InvalidRequest(err)
	}

	if err := h.brandService.Update(r.Context(), req, ID); err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Message: "Brand updated"})
}

func (h *BrandHandler) delete(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return util.InvalidParamID("id")
	}

	if err := h.brandService.Delete(r.Context(), ID); err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Message: "Brand deleted"})
}
