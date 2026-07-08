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

	roleRepo := repository.NewRoleRepository(s.db)
	userRepo := repository.NewUserRepository(s.db)
	authService := service.NewAuthService(userRepo, roleRepo)
	authHandler := handler.NewAuthHandler(authService)

	productRepo := repository.NewProductRepository(s.db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	cartRepo := repository.NewCartRepository(s.db)
	cartService := service.NewCartService(cartRepo, productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	r.Route("/auth", authHandler.Routes)
	r.Route("/products", productHandler.Routes)
	r.Route("/cart", cartHandler.Routes)

	log.Println("Server running on port", s.addr)
	return http.ListenAndServe(s.addr, r)
}
