// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client, { publicClient } from './client'

export interface Project {
  id: number
  name: string
  created_at: string
}

export interface App {
  id: number
  project_id: number
  name: string
  created_at: string
}

export interface Version {
  id: number
  app_id: number
  version: string
  created_at: string
}

export interface ProjectsResponse {
  projects: Project[]
  total: number
}

export const projectsApi = {
  list: (limit = 50, offset = 0) =>
    client.get<Project[]>('/projects', { params: { limit, offset } }),
  
  getApps: (project: string, limit = 50, offset = 0) =>
    client.get<App[]>(`/projects/${project}/apps`, { params: { limit, offset } }),
  
  getVersions: (project: string, app: string, limit = 50, offset = 0) =>
    client.get<Version[]>(`/projects/${project}/apps/${app}/versions`, {
      params: { limit, offset },
    }),
  
  deleteProject: (project: string) =>
    client.delete(`/projects/${project}`),
  
  deleteApp: (project: string, app: string) =>
    client.delete(`/projects/${project}/apps/${app}`),
  
  deleteVersion: (project: string, app: string, version: string) =>
    client.delete(`/projects/${project}/apps/${app}/versions/${version}`),
}

// Public API (no authentication required)
export const publicProjectsApi = {
  list: (limit = 50, offset = 0) =>
    publicClient.get<Project[]>('/public/projects', { params: { limit, offset } }),
  
  getApps: (project: string, limit = 50, offset = 0) =>
    publicClient.get<App[]>(`/public/projects/${project}/apps`, { params: { limit, offset } }),
  
  getVersions: (project: string, app: string, limit = 50, offset = 0) =>
    publicClient.get<Version[]>(`/public/projects/${project}/apps/${app}/versions`, {
      params: { limit, offset },
    }),
}

