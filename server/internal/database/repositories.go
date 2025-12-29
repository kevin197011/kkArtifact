// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"database/sql"
	"fmt"
)

// ProjectRepository handles project database operations
type ProjectRepository struct {
	db *DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// CreateOrGet creates a project or returns existing one
func (r *ProjectRepository) CreateOrGet(name string) (*Project, error) {
	var project Project
	query := `INSERT INTO projects (name) VALUES ($1) 
	          ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
	          RETURNING id, name, created_at`
	err := r.db.QueryRow(query, name).Scan(&project.ID, &project.Name, &project.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create or get project: %w", err)
	}
	return &project, nil
}

// List lists all projects ordered by creation time
func (r *ProjectRepository) List(limit, offset int) ([]*Project, error) {
	query := `SELECT id, name, created_at FROM projects 
	          ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, rows.Err()
}

// AppRepository handles app database operations
type AppRepository struct {
	db *DB
}

// NewAppRepository creates a new app repository
func NewAppRepository(db *DB) *AppRepository {
	return &AppRepository{db: db}
}

// CreateOrGet creates an app or returns existing one
func (r *AppRepository) CreateOrGet(projectID int, name string) (*App, error) {
	var app App
	query := `INSERT INTO apps (project_id, name) VALUES ($1, $2)
	          ON CONFLICT (project_id, name) DO UPDATE SET name = EXCLUDED.name
	          RETURNING id, project_id, name, created_at`
	err := r.db.QueryRow(query, projectID, name).Scan(&app.ID, &app.ProjectID, &app.Name, &app.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create or get app: %w", err)
	}
	return &app, nil
}

// ListByProject lists apps for a project
func (r *AppRepository) ListByProject(projectID int, limit, offset int) ([]*App, error) {
	query := `SELECT id, project_id, name, created_at FROM apps 
	          WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, projectID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	defer rows.Close()

	var apps []*App
	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.Name, &a.CreatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, &a)
	}
	return apps, rows.Err()
}

// VersionRepository handles version database operations
type VersionRepository struct {
	db *DB
}

// NewVersionRepository creates a new version repository
func NewVersionRepository(db *DB) *VersionRepository {
	return &VersionRepository{db: db}
}

// Create creates a new version
func (r *VersionRepository) Create(appID int, hash string) (*Version, error) {
	var version Version
	query := `INSERT INTO versions (app_id, hash) VALUES ($1, $2)
	          ON CONFLICT (app_id, hash) DO UPDATE SET app_id = EXCLUDED.app_id
	          RETURNING id, app_id, hash, created_at`
	err := r.db.QueryRow(query, appID, hash).Scan(&version.ID, &version.AppID, &version.Hash, &version.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}
	return &version, nil
}

// CreateOrGet creates a version or returns existing one (idempotent)
func (r *VersionRepository) CreateOrGet(appID int, hash string) (*Version, error) {
	var version Version
	query := `INSERT INTO versions (app_id, hash) VALUES ($1, $2)
	          ON CONFLICT (app_id, hash) DO NOTHING
	          RETURNING id, app_id, hash, created_at`
	err := r.db.QueryRow(query, appID, hash).Scan(&version.ID, &version.AppID, &version.Hash, &version.CreatedAt)
	if err != nil {
		// If no row returned (conflict), fetch existing version
		if err == sql.ErrNoRows {
			query := `SELECT id, app_id, hash, created_at FROM versions WHERE app_id = $1 AND hash = $2`
			err := r.db.QueryRow(query, appID, hash).Scan(&version.ID, &version.AppID, &version.Hash, &version.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to get existing version: %w", err)
			}
			return &version, nil
		}
		return nil, fmt.Errorf("failed to create or get version: %w", err)
	}
	return &version, nil
}

// ListByApp lists versions for an app
func (r *VersionRepository) ListByApp(appID int, limit, offset int) ([]*Version, error) {
	query := `SELECT id, app_id, hash, created_at FROM versions 
	          WHERE app_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, appID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}
	defer rows.Close()

	var versions []*Version
	for rows.Next() {
		var v Version
		if err := rows.Scan(&v.ID, &v.AppID, &v.Hash, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, &v)
	}
	return versions, rows.Err()
}

// Delete deletes a version
func (r *VersionRepository) Delete(appID int, hash string) error {
	query := `DELETE FROM versions WHERE app_id = $1 AND hash = $2`
	_, err := r.db.Exec(query, appID, hash)
	return err
}

// CountByApp counts versions for an app
func (r *VersionRepository) CountByApp(appID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM versions WHERE app_id = $1`
	err := r.db.QueryRow(query, appID).Scan(&count)
	return count, err
}

// GetOldestVersions returns oldest versions beyond limit
func (r *VersionRepository) GetOldestVersions(appID int, limit int) ([]*Version, error) {
	query := `SELECT id, app_id, hash, created_at FROM versions 
	          WHERE app_id = $1 
	          ORDER BY created_at ASC 
	          LIMIT $2`
	rows, err := r.db.Query(query, appID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest versions: %w", err)
	}
	defer rows.Close()

	var versions []*Version
	for rows.Next() {
		var v Version
		if err := rows.Scan(&v.ID, &v.AppID, &v.Hash, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, &v)
	}
	return versions, rows.Err()
}

