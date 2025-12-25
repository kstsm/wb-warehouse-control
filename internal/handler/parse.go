package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
)

func parseUUIDParam(r *http.Request) (uuid.UUID, error) {
	const param = "id"
	value := chi.URLParam(r, param)
	if strings.TrimSpace(value) == "" {
		return uuid.Nil, fmt.Errorf("%s is required", param)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", param)
	}

	return id, nil
}

func parseHistoryQuery(r *http.Request, req *dto.GetHistoryRequest) error {
	q := r.URL.Query()

	itemIDStr := strings.TrimSpace(q.Get("item_id"))
	if itemIDStr != "" {
		req.ItemID = &itemIDStr
	}

	userIDStr := strings.TrimSpace(q.Get("user_id"))
	if userIDStr != "" {
		req.UserID = &userIDStr
	}

	actionStr := strings.TrimSpace(q.Get("action"))
	if actionStr != "" {
		req.Action = &actionStr
	}

	if fromStr := q.Get("from"); fromStr != "" {
		from, err := parseDate(fromStr)
		if err != nil && !errors.Is(err, apperrors.ErrEmptyDate) {
			return err
		}
		if from != nil {
			req.From = from
		}
	}

	if toStr := q.Get("to"); toStr != "" {
		to, err := parseDate(toStr)
		if err != nil && !errors.Is(err, apperrors.ErrEmptyDate) {
			return err
		}
		if to != nil {
			req.To = to
		}
	}

	if req.From != nil && req.To != nil && req.From.After(*req.To) {
		return errors.New("parameter 'from' cannot be after 'to'")
	}

	sortByStr := strings.TrimSpace(q.Get("sort_by"))
	if sortByStr != "" {
		req.SortBy = &sortByStr
	}

	sortOrderStr := strings.TrimSpace(q.Get("sort_order"))
	if sortOrderStr != "" {
		req.SortOrder = &sortOrderStr
	}

	return nil
}

func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, apperrors.ErrEmptyDate
	}

	s = strings.TrimSpace(s)
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
