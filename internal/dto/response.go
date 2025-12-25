package dto

type ItemResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Price       string `json:"price"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ItemsListResponse struct {
	Items []ItemResponse `json:"items"`
	Total int            `json:"total"`
}

type HistoryResponse struct {
	ID        string         `json:"id"`
	ItemID    string         `json:"item_id"`
	Action    string         `json:"action"`
	UserID    *string        `json:"user_id,omitempty"`
	ChangedAt string         `json:"changed_at"`
	OldData   map[string]any `json:"old_data,omitempty"`
	NewData   map[string]any `json:"new_data,omitempty"`
}

type HistoryListResponse struct {
	History []HistoryResponse `json:"history"`
	Total   int               `json:"total"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type DiffResponse struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

type HistoryWithDiffResponse struct {
	HistoryResponse

	Diff []DiffResponse `json:"diff,omitempty"`
}

type HistoryWithDiffListResponse struct {
	History []HistoryWithDiffResponse `json:"history"`
	Total   int                       `json:"total"`
}

type HistoryExportResponse struct {
	ID        string `json:"id"`
	ItemID    string `json:"item_id"`
	Action    string `json:"action"`
	UserID    string `json:"user_id"`
	ChangedAt string `json:"changed_at"`
	OldData   string `json:"old_data"`
	NewData   string `json:"new_data"`
}

type ItemWithMessageResponse struct {
	ItemResponse

	Message string `json:"message,omitempty"`
}
