package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/converter"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
)

func (h *Handler) validateUUIDParams(req *dto.GetHistoryRequest) error {
	if req.ItemID != nil {
		if _, err := uuid.Parse(*req.ItemID); err != nil {
			return fmt.Errorf("invalid item_id: %w", err)
		}
	}

	if req.UserID != nil {
		if _, err := uuid.Parse(*req.UserID); err != nil {
			return fmt.Errorf("invalid user_id: %w", err)
		}
	}

	return nil
}

func (h *Handler) getItemHistoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetHistoryByItemID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "item not found")
		default:
			h.log.Errorf("Service error: %v", err)
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := converter.HistoriesToResponseWithDiff(result)

	h.respondJSON(w, http.StatusOK, dto.HistoryWithDiffListResponse{
		History: resp,
		Total:   len(resp),
	})
}

func (h *Handler) getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetHistoryRequest

	if err := parseHistoryQuery(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.valid.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, h.valid.FormatValidationError(err))
		return
	}

	if err := h.validateUUIDParams(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, total, err := h.service.GetHistory(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "item not found")
		default:
			h.log.Errorf("Service error: %v", err)
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := converter.HistoriesToResponse(result)
	h.respondJSON(w, http.StatusOK, dto.HistoryListResponse{
		History: resp,
		Total:   total,
	})
}

func (h *Handler) exportHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetHistoryRequest

	if err := parseHistoryQuery(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.valid.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, h.valid.FormatValidationError(err))
		return
	}

	if err := h.validateUUIDParams(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.ExportHistoryCSV(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "item not found")
		default:
			h.log.Errorf("Service error: %v", err)
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.respondCSV(w, http.StatusOK, data)
}
