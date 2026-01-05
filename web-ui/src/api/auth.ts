// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import client from './client'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  name: string
}

export const authApi = {
  login: (data: LoginRequest) => client.post<LoginResponse>('/login', data),
}

