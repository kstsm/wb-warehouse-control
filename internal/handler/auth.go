package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
)

func (h *Handler) SignInOrSignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.valid.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, h.valid.FormatValidationError(err))
		return
	}

	token, role, err := h.service.SignInOrSignUp(r.Context(), req.UserName, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserAlreadyExists):
			h.respondError(w, http.StatusConflict, "user already exists")
		case errors.Is(err, apperrors.ErrRoleMismatch):
			h.respondError(w, http.StatusConflict, "user already exists with role")
		default:
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.respondJSON(w, http.StatusOK, dto.LoginResponse{
		Token: token,
		Role:  role,
	})
}
