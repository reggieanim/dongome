-- Initial database setup for Dongome marketplace
-- This script runs when the PostgreSQL container is first created

-- Create extensions if they don't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- For text search

-- Create indexes for text search
-- These will be created by GORM auto-migration, but we can add custom ones here

-- Set up initial data
-- Categories will be populated via API or migrations