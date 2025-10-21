package database

import (
	"database/sql"
	"fmt"
)

type PackSizeRepository struct {
	db *DB
}

func NewPackSizeRepository(db *DB) *PackSizeRepository {
	return &PackSizeRepository{db: db}
}

// GetAllActive returns all active pack sizes sorted in ascending order
func (r *PackSizeRepository) GetAllActive() ([]int, error) {
	query := `SELECT size FROM pack_sizes WHERE is_active = true ORDER BY size ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active pack sizes: %w", err)
	}
	defer rows.Close()

	var sizes []int
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("failed to scan pack size: %w", err)
		}
		sizes = append(sizes, size)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pack sizes: %w", err)
	}

	return sizes, nil
}

// GetAll returns all pack sizes (active and inactive)
func (r *PackSizeRepository) GetAll() ([]PackSize, error) {
	query := `SELECT id, size, is_active, created_at, updated_at FROM pack_sizes ORDER BY size ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pack sizes: %w", err)
	}
	defer rows.Close()

	var packSizes []PackSize
	for rows.Next() {
		var ps PackSize
		if err := rows.Scan(&ps.ID, &ps.Size, &ps.IsActive, &ps.CreatedAt, &ps.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pack size: %w", err)
		}
		packSizes = append(packSizes, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pack sizes: %w", err)
	}

	return packSizes, nil
}

// GetByID returns a pack size by ID
func (r *PackSizeRepository) GetByID(id int) (*PackSize, error) {
	query := `SELECT id, size, is_active, created_at, updated_at FROM pack_sizes WHERE id = $1`
	
	var ps PackSize
	err := r.db.QueryRow(query, id).Scan(&ps.ID, &ps.Size, &ps.IsActive, &ps.CreatedAt, &ps.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pack size with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get pack size: %w", err)
	}

	return &ps, nil
}

// Create creates a new pack size
func (r *PackSizeRepository) Create(size int, isActive bool) (*PackSize, error) {
	query := `INSERT INTO pack_sizes (size, is_active) VALUES ($1, $2) RETURNING id, size, is_active, created_at, updated_at`
	
	var ps PackSize
	err := r.db.QueryRow(query, size, isActive).Scan(&ps.ID, &ps.Size, &ps.IsActive, &ps.CreatedAt, &ps.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create pack size: %w", err)
	}

	return &ps, nil
}

// Update updates an existing pack size
func (r *PackSizeRepository) Update(id int, size int, isActive bool) (*PackSize, error) {
	query := `UPDATE pack_sizes SET size = $1, is_active = $2 WHERE id = $3 RETURNING id, size, is_active, created_at, updated_at`
	
	var ps PackSize
	err := r.db.QueryRow(query, size, isActive, id).Scan(&ps.ID, &ps.Size, &ps.IsActive, &ps.CreatedAt, &ps.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pack size with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to update pack size: %w", err)
	}

	return &ps, nil
}

// Delete deletes a pack size (soft delete by setting is_active to false)
func (r *PackSizeRepository) Delete(id int) error {
	query := `UPDATE pack_sizes SET is_active = false WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pack size: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pack size with id %d not found", id)
	}

	return nil
}

// MigrateFromConfig migrates pack sizes from config to database
func (r *PackSizeRepository) MigrateFromConfig(sizes []int) error {
	// Get existing sizes from database
	existingSizes, err := r.GetAllActive()
	if err != nil {
		return fmt.Errorf("failed to get existing sizes: %w", err)
	}

	// Create a map for quick lookup
	existingMap := make(map[int]bool)
	for _, size := range existingSizes {
		existingMap[size] = true
	}

	// Insert new sizes that don't exist
	for _, size := range sizes {
		if !existingMap[size] {
			_, err := r.Create(size, true)
			if err != nil {
				return fmt.Errorf("failed to migrate size %d: %w", size, err)
			}
		}
	}

	return nil
}
