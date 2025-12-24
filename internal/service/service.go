package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/internal/repository"
	"github.com/kstsm/wb-warehouse-control/pkg/jwt"
)

type ItemManager interface {
	SignInOrSignUp(ctx context.Context, userName, role string) (string, string, error)
	CreateItem(ctx context.Context, req dto.CreateItemRequest, userID *uuid.UUID) (*models.Item, error)
	GetItems(ctx context.Context) ([]*models.Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequest, userID *uuid.UUID) (*models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error
	GetHistory(ctx context.Context, req dto.GetHistoryRequest) ([]*models.History, int, error)
	GetHistoryByItemID(ctx context.Context, itemID uuid.UUID) ([]*models.History, error)
	ExportHistoryCSV(ctx context.Context, req dto.GetHistoryRequest) ([]byte, error)
}

type Service struct {
	repo           repository.ItemManager
	log            *slog.Logger
	tokenGenerator jwt.TokenGenerator
}

func NewService(repo repository.ItemManager, log *slog.Logger, tokenGenerator jwt.TokenGenerator) ItemManager {
	return &Service{
		repo:           repo,
		log:            log,
		tokenGenerator: tokenGenerator,
	}
}
