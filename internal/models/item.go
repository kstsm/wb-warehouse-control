package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	Quantity    int
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
