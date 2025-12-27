// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface Manifest {
  project: string
  app: string
  version: string
  hash: string
  git_commit?: string
  build_time: string
  builder: string
  files: Array<{
    path: string
    size: number
    hash: string
  }>
}

export interface PromoteRequest {
  project: string
  app: string
  hash: string
}

export const versionsApi = {
  getManifest: (project: string, app: string, hash: string) =>
    client.get<Manifest>(`/manifest/${project}/${app}/${hash}`),
  promote: (data: PromoteRequest) => client.post('/promote', data),
  downloadFile: (project: string, app: string, hash: string, path: string) =>
    client.get(`/file/${project}/${app}/${hash}`, { params: { path }, responseType: 'blob' }),
}

