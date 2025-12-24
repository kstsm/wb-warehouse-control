package service

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kstsm/wb-warehouse-control/internal/converter"
	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/pkg/export"
)

func (s *Service) CreateItem(ctx context.Context, req dto.CreateItemRequest, userID *uuid.UUID) (*models.Item, error) {
	item := models.Item{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Price:       req.Price,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.CreateItem(ctx, item, userID); err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *Service) GetItems(ctx context.Context) ([]*models.Item, error) {
	return s.repo.GetItems(ctx)
}

func (s *Service) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	return s.repo.GetItemByID(ctx, id)
}

func (s *Service) UpdateItem(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateItemRequest,
	userID *uuid.UUID,
) (*models.Item, error) {
	return s.repo.UpdateItem(ctx, id, req, userID)
}

func (s *Service) DeleteItem(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error {
	return s.repo.DeleteItem(ctx, id, userID)
}

func (s *Service) GetHistory(ctx context.Context, req dto.GetHistoryRequest) ([]*models.History, int, error) {
	return s.repo.GetHistory(ctx, req)
}

func (s *Service) GetHistoryByItemID(ctx context.Context, itemID uuid.UUID) ([]*models.History, error) {
	return s.repo.GetHistoryByItemID(ctx, itemID)
}

func (s *Service) ExportHistoryCSV(ctx context.Context, req dto.GetHistoryRequest) ([]byte, error) {
	histories, _, err := s.GetHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	exportHistories := converter.HistoriesToExportResponse(histories)

	var buf bytes.Buffer
	if writeErr := export.WriteItemsCSV(&buf, nil, exportHistories); writeErr != nil {
		return nil, writeErr
	}

	return buf.Bytes(), nil
}
