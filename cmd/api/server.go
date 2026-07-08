package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonan98/ecom-pc-api/internal/handler"
	"github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr string
	db   *pgxpool.Pool
}

func NewServer(addr string, db *pgxpool.Pool) *Server {
	return &Server{
		addr: fmt.Sprintf(":%s", addr),
		db:   db,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cartRepo := repository.NewCartRepository(s.db)
	roleRepo := repository.NewRoleRepository(s.db)
	userRepo := repository.NewUserRepository(s.db)
	brandRepo := repository.NewBrandRepository(s.db)
	categoryRepo := repository.NewCategoryRepository(s.db)
	productRepo := repository.NewProductRepository(s.db)

	authService := service.NewAuthService(userRepo, roleRepo, cartRepo)
	brandService := service.NewBrandService(brandRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)

	authHandler := handler.NewAuthHandler(authService)
	brandHandler := handler.NewBrandHandler(brandService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", authHandler.Routes)
		r.Route("/brands", brandHandler.Routes)
		r.Route("/categories", categoryHandler.Routes)
		r.Route("/products", productHandler.Routes)
		r.Route("/cart", cartHandler.Routes)
	})

	log.Println("Server running on port", s.addr)
	return http.ListenAndServe(s.addr, r)
}
