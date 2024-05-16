package routes

import (
	"wbtech/level0/internal/api/handlers"
	"wbtech/level0/internal/api/middleware"
	"wbtech/level0/internal/config"
	"wbtech/level0/internal/customerrors"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux, h *handlers.Handlers, cfg *config.Config) {

	r.NotFound(customerrors.NotFoundResponse)
	r.MethodNotAllowed(customerrors.MethodNotAllowedResponse)

	r.Use(middleware.RecoverPanic)

	r.Get("/", h.IndexHandler.Handle)

	r.Route("/api/v1/", func(r chi.Router) {
		r.Mount("/orders", orderRoutes(h.OrderHandler))
	})
}

func orderRoutes(h *handlers.OrderHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/", h.Get)
	return r
}
