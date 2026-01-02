-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- Drop index
DROP INDEX IF EXISTS idx_versions_is_published;

-- Remove is_published column
ALTER TABLE versions DROP COLUMN IF EXISTS is_published;

