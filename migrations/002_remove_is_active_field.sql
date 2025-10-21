-- Migration: Remove is_active field since we're using hard deletes
-- Created: 2024-01-01

-- Drop the index on is_active
DROP INDEX IF EXISTS idx_pack_sizes_active;

-- Remove the is_active column
ALTER TABLE pack_sizes DROP COLUMN IF EXISTS is_active;

