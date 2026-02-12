import { api } from './api'
import type { UserSummary } from './types'

export async function listUsers(): Promise<Array<UserSummary>> {
  const response = await api.get<Array<UserSummary>>('/users')
  return response.data
}
