package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonan98/ecom-pc-api/internal/handler"
	"github.com/gonan98/ecom-pc-api/internal/service"
	"github.com/gonan98/ecom-pc-api/internal/store"
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

	roleStore := store.NewRoleStore(s.db)
	userStore := store.NewUserStore(s.db)
	authService := service.NewAuthService(userStore, roleStore)
	authHandler := handler.NewAuthHandler(authService)

	r.Route("/auth", authHandler.Routes)

	log.Println("Server running on port", s.addr)
	return http.ListenAndServe(s.addr, r)
}
