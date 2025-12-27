// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface Webhook {
  id: number
  name: string
  event_types: string[]
  url: string
  headers?: Record<string, string>
  enabled: boolean
  project_id?: number
  app_id?: number
  created_at: string
}

export interface CreateWebhookRequest {
  name: string
  event_types: string[]
  url: string
  headers?: Record<string, string>
  enabled?: boolean
  project_id?: number
  app_id?: number
}

export const webhooksApi = {
  list: () => client.get<Webhook[]>('/webhooks'),
  get: (id: number) => client.get<Webhook>(`/webhooks/${id}`),
  create: (data: CreateWebhookRequest) => client.post<Webhook>('/webhooks', data),
  update: (id: number, data: Partial<CreateWebhookRequest>) =>
    client.put<Webhook>(`/webhooks/${id}`, data),
  delete: (id: number) => client.delete(`/webhooks/${id}`),
}

