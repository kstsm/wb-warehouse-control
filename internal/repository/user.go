package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-warehouse-control/internal/apperrors"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/internal/repository/queries"
)

func (r *Repository) GetOrCreateUser(ctx context.Context, user models.User) (*models.User, error) {
	var resp models.User

	err := r.conn.QueryRow(ctx, queries.GetOrCreateUserQuery,
		user.ID,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&resp.ID,
		&resp.Name,
		&resp.Role,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("QueryRow-GetOrCreateUser: unexpected: no user returned")
		}
		return nil, fmt.Errorf("QueryRow-GetOrCreateUser: %w", err)
	}

	if resp.Role != user.Role {
		return nil, apperrors.ErrRoleMismatch
	}

	return &resp, nil
}

func (r *Repository) GetUserByName(ctx context.Context, userName string) (*models.User, error) {
	user := new(models.User)
	err := r.conn.QueryRow(ctx, queries.GetUserByNameQuery, userName).Scan(
		&user.ID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("QueryRow-GetUserByName: %w", err)
	}

	return user, nil
}
