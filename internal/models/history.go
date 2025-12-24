package models

import (
	"time"

	"github.com/google/uuid"
)

type History struct {
	ID        uuid.UUID
	ItemID    uuid.UUID
	Action    string
	UserID    *uuid.UUID
	ChangedAt time.Time
	OldData   map[string]any
	NewData   map[string]any
}
