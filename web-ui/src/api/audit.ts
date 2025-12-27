// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface AuditLog {
  id: number
  operation: string
  project_id?: number
  app_id?: number
  project_name?: string
  app_name?: string
  version_hash?: string
  agent_id?: string
  metadata?: string
  created_at: string
}

export const auditApi = {
  list: (limit = 50, offset = 0) =>
    client.get<AuditLog[]>('/audit-logs', { params: { limit, offset } }),
}

