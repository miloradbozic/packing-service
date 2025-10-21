package models

type CalculateRequest struct {
	Items int `json:"items"`
}

type CalculateResponse struct {
	Items       int    `json:"items_ordered"`
	TotalItems  int    `json:"total_items_shipped"`
	TotalPacks  int    `json:"total_packs"`
	Packs       []Pack `json:"packs"`
	ExcessItems int    `json:"excess_items"`
}

type Pack struct {
	Size     int `json:"size"`
	Quantity int `json:"quantity"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ConfigResponse struct {
	PackSizes []int `json:"pack_sizes"`
}

// Pack size management models
type PackSizeListResponse struct {
	PackSizes []PackSizeResponse `json:"pack_sizes"`
}

type PackSizeResponse struct {
	ID        int    `json:"id"`
	Size      int    `json:"size"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreatePackSizeRequest struct {
	Size     int  `json:"size"`
	IsActive bool `json:"is_active"`
}

type UpdatePackSizeRequest struct {
	Size     int  `json:"size"`
	IsActive bool `json:"is_active"`
}
