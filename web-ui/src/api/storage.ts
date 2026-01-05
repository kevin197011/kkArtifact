// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface SyncStorageResponse {
  message: string
  projects: number
  apps: number
  versions: number
}

export const storageApi = {
  syncStorage: () => client.post<SyncStorageResponse>('/sync-storage'),
}

