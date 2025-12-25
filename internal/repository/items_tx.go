package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/internal/repository/queries"
)

func setUserIDInTx(ctx context.Context, tx pgx.Tx, userID *uuid.UUID) error {
	if userID == nil {
		return nil
	}

	if _, err := tx.Exec(ctx, "SELECT set_config('app.user_id', $1, true)", userID.String()); err != nil {
		return fmt.Errorf("setUserIDInTx: %w", err)
	}

	return nil
}

func (r *Repository) CreateItem(ctx context.Context, item models.Item, userID *uuid.UUID) error {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("BeginTx-CreateItem: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Errorf("Rollback-CreateItem: %v", rbErr)
		}
	}()

	if errSetUser := setUserIDInTx(ctx, tx, userID); errSetUser != nil {
		return fmt.Errorf("setUserIDInTx-CreateItem: %w", errSetUser)
	}

	if _, err = tx.Exec(ctx, queries.CreateItemQuery,
		item.ID,
		item.Name,
		item.Description,
		item.Quantity,
		item.Price,
		item.CreatedAt,
		item.UpdatedAt,
	); err != nil {
		return fmt.Errorf("Exec-CreateItem: %w", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("Commit-CreateItem: %w", err)
	}

	return nil
}

func (r *Repository) UpdateItem(
	ctx context.Context,
	itemID uuid.UUID,
	req dto.UpdateItemRequest,
	userID *uuid.UUID,
) (*models.Item, error) {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("BeginTx-UpdateItem: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Errorf("Rollback-UpdateItem: %v", rbErr)
		}
	}()

	if errSetUser := setUserIDInTx(ctx, tx, userID); errSetUser != nil {
		return nil, fmt.Errorf("setUserIDInTx-UpdateItem: %w", errSetUser)
	}

	var item models.Item
	if err = tx.QueryRow(ctx, queries.UpdateItemQuery,
		itemID,
		req.Name,
		req.Description,
		req.Quantity,
		req.Price,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("QueryRow-UpdateItem: %w", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("Commit-UpdateItem: %w", err)
	}

	return &item, nil
}

func (r *Repository) DeleteItem(ctx context.Context, itemID uuid.UUID, userID *uuid.UUID) error {
	tx, err := r.conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("BeginTx-DeleteItem: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Errorf("Rollback-DeleteItem: %v", rbErr)
		}
	}()

	if errSetUser := setUserIDInTx(ctx, tx, userID); errSetUser != nil {
		return fmt.Errorf("setUserIDInTx-DeleteItem: %w", errSetUser)
	}

	var deletedID uuid.UUID
	if err = tx.QueryRow(ctx, queries.DeleteItemQuery, itemID).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrItemNotFound
		}
		return fmt.Errorf("QueryRow-DeleteItem: %w", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("Commit-DeleteItem: %w", err)
	}

	return nil
}
