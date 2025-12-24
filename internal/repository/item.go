package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/internal/repository/queries"
)

func (r *Repository) GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Item, error) {
	var item models.Item

	err := r.conn.QueryRow(ctx, queries.GetItemByIDQuery, itemID).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("QueryRow-GetItemByID: %w", err)
	}

	return &item, nil
}

func (r *Repository) GetItems(ctx context.Context) ([]*models.Item, error) {
	rows, err := r.conn.Query(ctx, queries.GetItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("Query-GetItems: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := new(models.Item)
		if errScan := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		); errScan != nil {
			return nil, fmt.Errorf("Scan-GetItems: %w", errScan)
		}
		items = append(items, item)
	}

	if errRows := rows.Err(); errRows != nil {
		return nil, fmt.Errorf("GetItems rows.Err: %w", errRows)
	}

	return items, nil
}

func (r *Repository) GetHistoryByItemID(ctx context.Context, itemID uuid.UUID) ([]*models.History, error) {
	rows, err := r.conn.Query(ctx, queries.GetHistoryByItemIDQuery, itemID)
	if err != nil {
		return nil, fmt.Errorf("Query-GetHistoryByItemID: %w", err)
	}
	defer rows.Close()

	histories, err := r.scanHistories(rows)
	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *Repository) GetHistory(ctx context.Context, req dto.GetHistoryRequest) ([]*models.History, int, error) {
	whereClause, args := r.buildHistoryWhere(req)
	orderClause := r.buildHistoryOrder(req)

	var total int
	countQuery := fmt.Sprintf(queries.GetHistoryCountQuery, whereClause)
	if err := r.conn.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("QueryRow-GetHistory: %w", err)
	}

	query := fmt.Sprintf(queries.GetHistoryQuery, whereClause, orderClause)
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Query-GetHistory: %w", err)
	}
	defer rows.Close()

	histories, err := r.scanHistories(rows)
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *Repository) scanHistories(rows pgx.Rows) ([]*models.History, error) {
	var histories []*models.History
	for rows.Next() {
		history := new(models.History)
		var oldDataJSON, newDataJSON []byte

		if err := rows.Scan(
			&history.ID,
			&history.ItemID,
			&history.Action,
			&history.UserID,
			&history.ChangedAt,
			&oldDataJSON,
			&newDataJSON,
		); err != nil {
			return nil, fmt.Errorf("scanHistories scan: %w", err)
		}

		if len(oldDataJSON) > 0 {
			if err := json.Unmarshal(oldDataJSON, &history.OldData); err != nil {
				return nil, fmt.Errorf("scanHistories unmarshal oldData: %w", err)
			}
		}

		if len(newDataJSON) > 0 {
			if err := json.Unmarshal(newDataJSON, &history.NewData); err != nil {
				return nil, fmt.Errorf("scanHistories unmarshal newData: %w", err)
			}
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanHistories rows.Err: %w", err)
	}

	return histories, nil
}

func (r *Repository) buildHistoryWhere(req dto.GetHistoryRequest) (string, []any) {
	var cond []string
	var args []any

	add := func(query string, val any) {
		cond = append(cond, fmt.Sprintf(query, len(args)+1))
		args = append(args, val)
	}

	if req.ItemID != nil {
		add("item_id = $%d::uuid", *req.ItemID)
	}
	if req.UserID != nil {
		add("user_id = $%d::uuid", *req.UserID)
	}
	if req.Action != nil {
		add("action = $%d", *req.Action)
	}
	if req.From != nil {
		add("changed_at >= $%d", *req.From)
	}
	if req.To != nil {
		add("changed_at <= $%d", *req.To)
	}

	if len(cond) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(cond, " AND "), args
}

func (r *Repository) buildHistoryOrder(req dto.GetHistoryRequest) string {
	sortBy := "changed_at"
	sortOrder := "DESC"

	if req.SortBy != nil {
		allowedSortBy := map[string]string{
			"changed_at": "changed_at",
			"action":     "action",
			"user_id":    "user_id",
		}
		if allowed, ok := allowedSortBy[strings.ToLower(*req.SortBy)]; ok {
			sortBy = allowed
		}
	}
	if req.SortOrder != nil {
		order := strings.ToUpper(*req.SortOrder)
		if order == "ASC" || order == "DESC" {
			sortOrder = order
		}
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}
