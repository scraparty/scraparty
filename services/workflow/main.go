package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Post("/workflow/new", func(w http.ResponseWriter, r *http.Request) {})

	r.Get("/workflow/:id", func(w http.ResponseWriter, r *http.Request) {})

	r.Get("/workflow/list", func(w http.ResponseWriter, r *http.Request) {})

	http.ListenAndServe(":3000", r)
}
