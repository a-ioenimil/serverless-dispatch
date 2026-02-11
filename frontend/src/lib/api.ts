import axios from 'axios'

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL as string | undefined

if (!apiBaseUrl) {
  console.warn('VITE_API_BASE_URL is not set; API calls will fail.')
}

export const api = axios.create({
  baseURL: apiBaseUrl,
  timeout: 8000,
})

api.interceptors.request.use((config) => {
  const token =
    localStorage.getItem('access_token') || localStorage.getItem('id_token')

  if (token) {
    config.headers = {
      ...config.headers,
      Authorization: `Bearer ${token}`,
    }
  }

  return config
})
