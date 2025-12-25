package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-warehouse-control/internal/middleware"
	"github.com/kstsm/wb-warehouse-control/internal/service"
	"github.com/kstsm/wb-warehouse-control/pkg/jwt"
	"github.com/kstsm/wb-warehouse-control/pkg/validator"
)

type ItemManager interface {
	NewRouter() http.Handler
}

type Handler struct {
	service        service.ItemManager
	log            *slog.Logger
	valid          *validator.Validate
	tokenValidator jwt.TokenValidator
}

func NewHandler(
	service service.ItemManager,
	log *slog.Logger,
	valid *validator.Validate,
	tokenValidator jwt.TokenValidator,
) ItemManager {
	return &Handler{
		service:        service,
		log:            log,
		valid:          valid,
		tokenValidator: tokenValidator,
	}
}

func (h *Handler) NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS)

	h.registerPublicRoutes(r)
	h.registerAPIRoutes(r)
	h.registerWebRoutes(r)

	return r
}

func (h *Handler) serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/"+filename)
	}
}
