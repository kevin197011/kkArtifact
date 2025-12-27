// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface Config {
  version_retention_limit: number
}

export const configApi = {
  get: () => client.get<Config>('/config'),
  update: (data: Partial<Config>) => client.put('/config', data),
}

