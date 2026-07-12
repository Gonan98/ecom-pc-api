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
	r.Get("/", httpHandler(h.GetAll))
	r.Get("/{id}", httpHandler(h.GetByID))

	r.With(
		middleware.JWTMiddleware,
		middleware.AdminMiddleware,
	).Post("/", httpHandler(h.create))
}

func (h *BrandHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	brands, err := h.brandService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return write(w, types.APIResponse{Code: http.StatusOK, Data: brands})
}

func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return types.NewAPIError(http.StatusBadRequest, err)
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
