package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/stageddat/shelter-node/internal/db"
)

type Server struct {
	router chi.Router
	store  db.Store
}

func New(store db.Store) *Server {
	s := &Server{store: store}
	s.router = s.setupRouter()
	return s
}

func (s *Server) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://shelter.cat", "https://localhost:5173", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.registerRoutes(r)
	return r
}

func (s *Server) Start(addr string) error {
	fmt.Printf("starting server on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
