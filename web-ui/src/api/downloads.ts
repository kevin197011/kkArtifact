// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { publicClient } from './client'

export interface AgentBinaryInfo {
  platform: string
  filename: string
  size: number
  url: string
}

export interface AgentVersionInfo {
  version: string
  build_time: string
  binaries: AgentBinaryInfo[]
}

export const downloadsApi = {
  getAgentVersionInfo: () => publicClient.get<AgentVersionInfo>('/downloads/agent/version'),
  
  downloadAgent: (filename: string) => {
    const API_URL = import.meta.env.VITE_API_URL || '/'
    const baseURL = (API_URL === '/' || API_URL === '')
      ? '/api/v1'  // Relative path for nginx proxy
      : (API_URL.endsWith('/api/v1') 
          ? API_URL 
          : `${API_URL}/api/v1`)
    return `${baseURL}/downloads/agent/${filename}`
  },
}

