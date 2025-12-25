package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/pkg/jwt"
)

type contextKey string

const (
	userIDContextKey contextKey = "user_id"
	roleContextKey   contextKey = "role"
)

func AuthMiddleware(tokenValidator jwt.TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := extractBearerToken(r)
			if !ok {
				respondError(w)
				return
			}

			claims, err := tokenValidator.ValidateToken(token)
			if err != nil {
				respondError(w)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, roleContextKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowed ...jwt.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok {
				respondError(w)
				return
			}

			if slices.Contains(allowed, role) {
				next.ServeHTTP(w, r)
				return
			}

			respondError(w)
		})
	}
}

func RoleFromContext(ctx context.Context) (jwt.Role, bool) {
	role, ok := ctx.Value(roleContextKey).(jwt.Role)
	return role, ok
}

func UserIDFromContext(ctx context.Context) (*uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return nil, false
	}
	return &userID, true
}

func extractBearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	if !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}

func respondError(w http.ResponseWriter) {
	const (
		status  = http.StatusUnauthorized
		message = "unauthorized"
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(models.Error{Error: message})
	if err != nil {
		slog.Errorf("respondError: %v", err.Error())
		return
	}
}
