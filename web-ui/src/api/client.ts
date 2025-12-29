// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import axios from 'axios'

// Get API URL from environment variable
// For production nginx deployment, use relative path '/'
// For development, use full URL like 'http://localhost:8080'
const API_URL = import.meta.env.VITE_API_URL || '/'

// Use relative path if API_URL is '/' or empty
// Otherwise use the full path
const baseURL = (API_URL === '/' || API_URL === '')
  ? '/api/v1'  // Relative path for nginx proxy
  : (API_URL.endsWith('/api/v1') 
      ? API_URL 
      : `${API_URL}/api/v1`)

const client = axios.create({
  baseURL: baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Create a public client that doesn't add auth token or redirect on 401
export const publicClient = axios.create({
  baseURL: baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Public client should not redirect on 401 errors - it's expected for public endpoints
// No response interceptor needed for public client

// Add auth token interceptor
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('kkartifact_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle 401 errors - redirect to login
client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('kkartifact_token')
      // Only redirect if we're not already on login page
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default client

