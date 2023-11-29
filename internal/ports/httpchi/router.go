package httpchi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
)

func NewRouter(service Service, logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(LoggerRequestID(logger))
	r.Use(middleware.URLFormat)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)

	r.NotFound(http.NotFound)
	RegisterRoutes(r, service)
	return r
}

func RegisterRoutes(r *chi.Mux, service Service) {
	api := chi.NewRouter()

	api.Post("/task", service.CreateTask)
	api.Get("/tasks", nil)
	api.Get("/task/{id}", service.GetSingleTask)
	api.Put("/task/{id}", service.UpdateTask)
	api.Delete("/task/{id}", service.DeleteTask)

	r.Mount("/api", api)
}
