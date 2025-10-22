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


// GetAll returns all pack sizes
func (r *PackSizeRepository) GetAll() ([]PackSize, error) {
	query := `SELECT id, size, created_at, updated_at FROM pack_sizes ORDER BY size ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pack sizes: %w", err)
	}
	defer rows.Close()

	var packSizes []PackSize
	for rows.Next() {
		var ps PackSize
		if err := rows.Scan(&ps.ID, &ps.Size, &ps.CreatedAt, &ps.UpdatedAt); err != nil {
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
	query := `SELECT id, size, created_at, updated_at FROM pack_sizes WHERE id = $1`
	
	var ps PackSize
	err := r.db.QueryRow(query, id).Scan(&ps.ID, &ps.Size, &ps.CreatedAt, &ps.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pack size with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get pack size: %w", err)
	}

	return &ps, nil
}

// Create creates a new pack size
func (r *PackSizeRepository) Create(size int) (*PackSize, error) {
	query := `INSERT INTO pack_sizes (size) VALUES ($1) RETURNING id, size, created_at, updated_at`
	
	var ps PackSize
	err := r.db.QueryRow(query, size).Scan(&ps.ID, &ps.Size, &ps.CreatedAt, &ps.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create pack size: %w", err)
	}

	return &ps, nil
}

// Update updates an existing pack size
func (r *PackSizeRepository) Update(id int, size int) (*PackSize, error) {
	query := `UPDATE pack_sizes SET size = $1 WHERE id = $2 RETURNING id, size, created_at, updated_at`
	
	var ps PackSize
	err := r.db.QueryRow(query, size, id).Scan(&ps.ID, &ps.Size, &ps.CreatedAt, &ps.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pack size with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to update pack size: %w", err)
	}

	return &ps, nil
}

// Delete deletes a pack size (hard delete - removes from database)
func (r *PackSizeRepository) Delete(id int) error {
	query := `DELETE FROM pack_sizes WHERE id = $1`
	
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

