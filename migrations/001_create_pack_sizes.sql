-- Migration: Create pack_sizes table
-- Created: 2024-01-01

CREATE TABLE IF NOT EXISTS pack_sizes (
    id SERIAL PRIMARY KEY,
    size INTEGER NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on size for faster lookups
CREATE INDEX IF NOT EXISTS idx_pack_sizes_size ON pack_sizes(size);
CREATE INDEX IF NOT EXISTS idx_pack_sizes_active ON pack_sizes(is_active);

-- Insert default pack sizes
INSERT INTO pack_sizes (size, is_active) VALUES 
    (250, true),
    (500, true),
    (1000, true),
    (2000, true),
    (5000, true)
ON CONFLICT (size) DO NOTHING;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_pack_sizes_updated_at ON pack_sizes;
CREATE TRIGGER update_pack_sizes_updated_at 
    BEFORE UPDATE ON pack_sizes 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
