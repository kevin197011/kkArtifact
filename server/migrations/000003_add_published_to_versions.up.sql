-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- Add is_published column to versions table
ALTER TABLE versions ADD COLUMN IF NOT EXISTS is_published BOOLEAN NOT NULL DEFAULT FALSE;

-- Create index for faster queries of published versions
CREATE INDEX IF NOT EXISTS idx_versions_is_published ON versions(app_id, is_published, created_at DESC) WHERE is_published = TRUE;

