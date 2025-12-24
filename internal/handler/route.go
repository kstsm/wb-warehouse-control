package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/kstsm/wb-warehouse-control/internal/middleware"
	"github.com/kstsm/wb-warehouse-control/pkg/jwt"
)

func (h *Handler) registerPublicRoutes(r chi.Router) {
	r.Post("/api/login", h.SignInOrSignUp)
}

func (h *Handler) registerAPIRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(h.tokenValidator))

		r.Route("/items", func(r chi.Router) {
			r.With(middleware.RequireRole(jwt.RoleAdmin, jwt.RoleManager)).Group(func(r chi.Router) {
				r.Post("/", h.createItemHandler)
				r.Put("/{id}", h.updateItemHandler)
			})

			r.With(middleware.RequireRole(jwt.RoleAdmin)).Delete("/{id}", h.deleteItemHandler)

			r.Get("/", h.getItemsHandler)
			r.Get("/{id}", h.getItemByIDHandler)
			r.Get("/{id}/history", h.getItemHistoryHandler)
		})

		r.Route("/history", func(r chi.Router) {
			r.Get("/", h.getHistoryHandler)
			r.Get("/export", h.exportHistoryHandler)
		})
	})
}

func (h *Handler) registerWebRoutes(r chi.Router) {
	r.Get("/", h.serveHTML("index.html"))
}
