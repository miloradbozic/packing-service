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
