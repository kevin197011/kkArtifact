// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface Token {
  id: number
  name: string
  token?: string // Only present when creating
  project_id?: number
  app_id?: number
  permissions: string[]
  expires_at?: string
  created_at: string
}

export interface CreateTokenRequest {
  name: string
  project_id?: number
  app_id?: number
  permissions?: string[]
  expires_at?: string
}

export const tokensApi = {
  list: () => client.get<Token[]>('/tokens'),
  create: (data: CreateTokenRequest) => client.post<Token>('/tokens', data),
  delete: (id: number) => client.delete(`/tokens/${id}`),
}

