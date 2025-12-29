-- Copyright (c) 2025 kk
--
-- This software is released under the MIT License.
-- https://opensource.org/licenses/MIT

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Apps table
CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, name)
);

-- Versions table
CREATE TABLE IF NOT EXISTS versions (
    id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(app_id, hash)
);

-- Tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
    permissions TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Webhooks table
CREATE TABLE IF NOT EXISTS webhooks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    event_types TEXT[] NOT NULL,
    url TEXT NOT NULL,
    headers JSONB,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    operation VARCHAR(50) NOT NULL,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    app_id INTEGER REFERENCES apps(id) ON DELETE SET NULL,
    version_hash VARCHAR(255),
    agent_id VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Global config table
CREATE TABLE IF NOT EXISTS config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_apps_project_id ON apps(project_id);
CREATE INDEX IF NOT EXISTS idx_versions_app_id ON versions(app_id);
CREATE INDEX IF NOT EXISTS idx_versions_created_at ON versions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tokens_project_id ON tokens(project_id);
CREATE INDEX IF NOT EXISTS idx_tokens_app_id ON tokens(app_id);
CREATE INDEX IF NOT EXISTS idx_webhooks_project_id ON webhooks(project_id);
CREATE INDEX IF NOT EXISTS idx_webhooks_app_id ON webhooks(app_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_project_id ON audit_logs(project_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_app_id ON audit_logs(app_id);

-- Insert default config
INSERT INTO config (key, value) VALUES ('version_retention_limit', '30')
ON CONFLICT (key) DO NOTHING;

