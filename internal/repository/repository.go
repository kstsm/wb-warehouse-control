package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/models"
)

type ItemManager interface {
	CreateItem(ctx context.Context, item models.Item, userID *uuid.UUID) error
	GetOrCreateUser(ctx context.Context, user models.User) (*models.User, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	GetItems(ctx context.Context) ([]*models.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequest, userID *uuid.UUID) (*models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error
	GetHistory(ctx context.Context, req dto.GetHistoryRequest) ([]*models.History, int, error)
	GetHistoryByItemID(ctx context.Context, itemID uuid.UUID) ([]*models.History, error)
}

type Repository struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func NewRepository(conn *pgxpool.Pool, log *slog.Logger) ItemManager {
	return &Repository{
		conn: conn,
		log:  log,
	}
}
