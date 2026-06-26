// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client, { publicClient } from './client'
import type { App, Project, Version } from './projects'

export interface AppInventory {
  app: App
  versions: Version[]
}

export interface ProjectInventory {
  project: Project
  apps: AppInventory[]
}

export interface InventoryResponse {
  projects: ProjectInventory[]
}

export interface InventorySummary {
  total_projects: number
  total_apps: number
  total_versions: number
}

export const publicInventoryApi = {
  getComplete: () => publicClient.get<InventoryResponse>('/public/inventory'),
}

export const inventoryApi = {
  getSummary: () => client.get<InventorySummary>('/admin/inventory/summary'),
}
