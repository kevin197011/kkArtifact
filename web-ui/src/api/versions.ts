// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface Manifest {
  project: string
  app: string
  version: string
  git_commit?: string
  build_time: string
  builder: string
  files: Array<{
    path: string
    size: number
    hash: string  // File SHA256 hash, not version identifier
  }>
}

export interface PromoteRequest {
  project: string
  app: string
  version: string
}

export const versionsApi = {
  // version parameter is the version identifier (stored as hash in database, but exposed as version in API)
  getManifest: (project: string, app: string, version: string) =>
    client.get<Manifest>(`/manifest/${project}/${app}/${version}`),
  promote: (data: PromoteRequest) => client.post('/promote', data),
  downloadFile: (project: string, app: string, version: string, path: string) =>
    client.get(`/file/${project}/${app}/${version}`, { params: { path }, responseType: 'blob' }),
}

