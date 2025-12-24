package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/converter"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/middleware"
)

func (h *Handler) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.valid.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, h.valid.FormatValidationError(err))
		return
	}

	var userID *uuid.UUID
	if id, ok := middleware.UserIDFromContext(r.Context()); ok {
		userID = id
	}

	result, err := h.service.CreateItem(r.Context(), req, userID)
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

	resp := converter.ItemToResponse(result)
	h.respondJSON(w, http.StatusCreated, dto.ItemWithMessageResponse{
		ItemResponse: resp,
		Message:     "item created successfully",
	})
}

func (h *Handler) getItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseUUIDParam(r, "id")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetItemByID(r.Context(), itemID)
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

	resp := converter.ItemToResponse(result)
	h.respondJSON(w, http.StatusOK, resp)
}

func (h *Handler) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetItems(r.Context())
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

	resp := converter.ItemsToResponse(result)
	h.respondJSON(w, http.StatusOK, dto.ItemsListResponse{
		Items: resp,
		Total: len(resp),
	})
}

func (h *Handler) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseUUIDParam(r, "id")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateItemRequest
	if errDecode := json.NewDecoder(r.Body).Decode(&req); errDecode != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errValidate := h.valid.Struct(req); errValidate != nil {
		h.respondError(w, http.StatusBadRequest, h.valid.FormatValidationError(errValidate))
		return
	}

	var userID *uuid.UUID
	if id, ok := middleware.UserIDFromContext(r.Context()); ok {
		userID = id
	}

	result, err := h.service.UpdateItem(r.Context(), itemID, req, userID)
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

	resp := converter.ItemToResponse(result)
	h.respondJSON(w, http.StatusOK, dto.ItemWithMessageResponse{
		ItemResponse: resp,
		Message:     "item updated successfully",
	})
}

func (h *Handler) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseUUIDParam(r, "id")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var userID *uuid.UUID
	if id, ok := middleware.UserIDFromContext(r.Context()); ok {
		userID = id
	}

	err = h.service.DeleteItem(r.Context(), itemID, userID)
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

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "item deleted successfully"})
}
