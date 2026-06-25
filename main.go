package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "shelter node is running :)")
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotImplemented)
				json.NewEncoder(w).Encode(struct {
					Status string `json:"status"`
				}{
					Status: "not implemented",
				})
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})
		})

		r.Route("/entries", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})
			r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})
			r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})
		})
	})

	fmt.Println("server running on port 4123 :)")
	http.ListenAndServe(":4123", r)
}
