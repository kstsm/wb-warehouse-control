package converter

import (
	"encoding/json"
	"reflect"
	"sort"
	"time"

	"github.com/kstsm/wb-warehouse-control/internal/dto"
	"github.com/kstsm/wb-warehouse-control/internal/models"
)

func HistoryToResponse(history *models.History) dto.HistoryResponse {
	var userID *string
	if history.UserID != nil {
		id := history.UserID.String()
		userID = &id
	}

	return dto.HistoryResponse{
		ID:        history.ID.String(),
		ItemID:    history.ItemID.String(),
		Action:    history.Action,
		UserID:    userID,
		ChangedAt: history.ChangedAt.UTC().Format(time.RFC3339),
		OldData:   history.OldData,
		NewData:   history.NewData,
	}
}

func HistoriesToResponse(histories []*models.History) []dto.HistoryResponse {
	res := make([]dto.HistoryResponse, len(histories))
	for i, h := range histories {
		res[i] = HistoryToResponse(h)
	}

	return res
}

func HistoryToResponseWithDiff(history *models.History) dto.HistoryWithDiffResponse {
	diff := calculateDiff(history.OldData, history.NewData)

	return dto.HistoryWithDiffResponse{
		HistoryResponse: HistoryToResponse(history),
		Diff:            diff,
	}
}

func HistoriesToResponseWithDiff(histories []*models.History) []dto.HistoryWithDiffResponse {
	res := make([]dto.HistoryWithDiffResponse, len(histories))
	for i, h := range histories {
		res[i] = HistoryToResponseWithDiff(h)
	}

	return res
}

func calculateDiff(oldData, newData map[string]any) []dto.DiffResponse {
	var diff []dto.DiffResponse

	allKeys := make(map[string]bool)
	for k := range oldData {
		allKeys[k] = true
	}
	for k := range newData {
		allKeys[k] = true
	}

	sortedKeys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		oldVal, oldExists := oldData[key]
		newVal, newExists := newData[key]

		switch {
		case !oldExists:
			diff = append(diff, dto.DiffResponse{
				Field:    key,
				OldValue: nil,
				NewValue: newVal,
			})
		case !newExists:
			diff = append(diff, dto.DiffResponse{
				Field:    key,
				OldValue: oldVal,
				NewValue: nil,
			})
		case !reflect.DeepEqual(oldVal, newVal):
			diff = append(diff, dto.DiffResponse{
				Field:    key,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	return diff
}

func HistoryToExportResponse(history *models.History) dto.HistoryExportResponse {
	var userID string
	if history.UserID != nil {
		userID = history.UserID.String()
	}

	oldDataStr := ""
	if history.OldData != nil {
		oldDataBytes, err := json.Marshal(history.OldData)
		if err == nil {
			oldDataStr = string(oldDataBytes)
		}
	}

	newDataStr := ""
	if history.NewData != nil {
		newDataBytes, err := json.Marshal(history.NewData)
		if err == nil {
			newDataStr = string(newDataBytes)
		}
	}

	return dto.HistoryExportResponse{
		ID:        history.ID.String(),
		ItemID:    history.ItemID.String(),
		Action:    history.Action,
		UserID:    userID,
		ChangedAt: history.ChangedAt.UTC().Format(time.RFC3339),
		OldData:   oldDataStr,
		NewData:   newDataStr,
	}
}

func HistoriesToExportResponse(histories []*models.History) []dto.HistoryExportResponse {
	res := make([]dto.HistoryExportResponse, len(histories))
	for i, h := range histories {
		res[i] = HistoryToExportResponse(h)
	}

	return res
}
