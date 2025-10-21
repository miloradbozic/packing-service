package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

type Migrator struct {
	db *DB
}

func NewMigrator(db *DB) *Migrator {
	return &Migrator{db: db}
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations(migrationsPath string) error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := m.getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migrationFile := range migrationFiles {
		if !m.isMigrationApplied(migrationFile, appliedMigrations) {
			if err := m.applyMigration(migrationFile, migrationsPath); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migrationFile, err)
			}
		}
	}

	return nil
}

func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getMigrationFiles(migrationsPath string) ([]string, error) {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func (m *Migrator) getAppliedMigrations() ([]string, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied = append(applied, version)
	}

	return applied, nil
}

func (m *Migrator) isMigrationApplied(migrationFile string, appliedMigrations []string) bool {
	version := strings.TrimSuffix(migrationFile, ".sql")
	for _, applied := range appliedMigrations {
		if applied == version {
			return true
		}
	}
	return false
}

func (m *Migrator) applyMigration(migrationFile, migrationsPath string) error {
	// Read migration file
	migrationPath := filepath.Join(migrationsPath, migrationFile)
	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	// Execute migration
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(content))
	if err != nil {
		return err
	}

	// Record migration as applied
	version := strings.TrimSuffix(migrationFile, ".sql")
	_, err = tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, version)
	if err != nil {
		return err
	}

	return tx.Commit()
}
