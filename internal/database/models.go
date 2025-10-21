package database

import (
	"time"
)

// PackSize represents a pack size configuration in the database
type PackSize struct {
	ID        int       `json:"id" db:"id"`
	Size      int       `json:"size" db:"size"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PackSizeRequest represents a request to create/update a pack size
type PackSizeRequest struct {
	Size int `json:"size"`
}

// PackSizeResponse represents the response for pack size operations
type PackSizeResponse struct {
	ID        int       `json:"id"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
