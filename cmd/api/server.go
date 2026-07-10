package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonan98/ecom-pc-api/internal/database"
	"github.com/gonan98/ecom-pc-api/internal/handler"
	repo "github.com/gonan98/ecom-pc-api/internal/repository"
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

	cartRepo := repo.NewCartRepository(s.db)
	roleRepo := repo.NewRoleRepository(s.db)
	userRepo := repo.NewUserRepository(s.db)
	brandRepo := repo.NewBrandRepository(s.db)
	categoryRepo := repo.NewCategoryRepository(s.db)
	productRepo := repo.NewProductRepository(s.db)
	orderRepo := repo.NewOrderRepository(s.db)

	txManager := database.NewTxManager(s.db)

	authService := service.NewAuthService(userRepo, roleRepo, cartRepo, txManager)
	brandService := service.NewBrandService(brandRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, cartRepo, txManager)

	authHandler := handler.NewAuthHandler(authService)
	brandHandler := handler.NewBrandHandler(brandService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", authHandler.Routes)
		r.Route("/brands", brandHandler.Routes)
		r.Route("/categories", categoryHandler.Routes)
		r.Route("/products", productHandler.Routes)
		r.Route("/cart", cartHandler.Routes)
		r.Route("/orders", orderHandler.Routes)
	})

	log.Println("Server running on port", s.addr)
	return http.ListenAndServe(s.addr, r)
}
