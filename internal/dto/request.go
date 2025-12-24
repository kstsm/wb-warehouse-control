package dto

import (
	"time"
)

type CreateItemRequest struct {
	Name        string `json:"name"        validate:"required,min=1"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"    validate:"required,min=1"`
	Price       int    `json:"price"       validate:"required,min=1"`
}

type UpdateItemRequest struct {
	Name        *string `json:"name"        validate:"omitempty,min=1"`
	Description *string `json:"description"`
	Quantity    *int    `json:"quantity"    validate:"omitempty,min=0"`
	Price       *int    `json:"price"       validate:"omitempty,min=0"`
}

type LoginRequest struct {
	UserName string `json:"user_name" validate:"required,letters_only"`
	Role     string `json:"role"      validate:"required,role"`
}

type GetHistoryRequest struct {
	ItemID    *string    `json:"item_id"`
	UserID    *string    `json:"user_id"`
	Action    *string    `json:"action"     validate:"omitempty,action_type"`
	From      *time.Time `json:"from"`
	To        *time.Time `json:"to"`
	SortBy    *string    `json:"sort_by"`
	SortOrder *string    `json:"sort_order"`
}
