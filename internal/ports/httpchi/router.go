package httpchi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vlasashk/todo-manager/internal/models/todo"
	"net/http"
)

func NewRouter(db todo.Repo) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)

	r.NotFound(http.NotFound)
	RegisterRoutes(r, db)
	return r
}

func RegisterRoutes(r *chi.Mux, db todo.Repo) {
	api := chi.NewRouter()

	api.Post("/task", nil)
	api.Get("/tasks", nil)
	api.Get("/task/{id}", nil)
	api.Put("/task/{id}", nil)
	api.Delete("/task/{id}", nil)

	r.Mount("/api", api)
}
