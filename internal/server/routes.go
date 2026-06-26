package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stageddat/shelter-node/internal/entry"
	"github.com/stageddat/shelter-node/internal/user"
)

func (s *Server) registerRoutes(r chi.Router) {
	userHandler := user.NewHandler(user.NewService(s.store))
	entryHandler := entry.NewHandler(entry.NewService(s.store))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("shelter node is running :)"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/user", userHandler.Routes)
		r.Route("/entries", entryHandler.Routes)
	})
}
