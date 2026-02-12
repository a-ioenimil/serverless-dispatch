import axios from 'axios'
import type { AxiosRequestHeaders } from 'axios'
import { getStoredTokens } from './auth'

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL as string | undefined

if (!apiBaseUrl) {
  console.warn('VITE_API_BASE_URL is not set; API calls will fail.')
}

export const api = axios.create({
  baseURL: apiBaseUrl,
  timeout: 8000,
})

api.interceptors.request.use((config) => {
  const tokens = getStoredTokens()

  if (tokens?.accessToken) {
    config.headers = {
      ...config.headers,
      Authorization: `Bearer ${tokens.accessToken}`,
    } as AxiosRequestHeaders
  }

  return config
})
