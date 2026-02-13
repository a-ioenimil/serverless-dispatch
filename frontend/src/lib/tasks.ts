import { api } from './api'
import type { CreateTaskInput, Task, UpdateTaskInput } from './types'

function parseApiBody<T>(data: unknown): T | null {
  if (typeof data === 'string') {
    try {
      return JSON.parse(data) as T
    } catch {
      return null
    }
  }

  if (data && typeof data === 'object' && 'body' in data) {
    const body = (data as { body?: unknown }).body
    if (typeof body === 'string') {
      try {
        return JSON.parse(body) as T
      } catch {
        return null
      }
    }
    return body as T
  }

  return data as T
}

export async function listTasks(): Promise<Array<Task>> {
  const response = await api.get<Array<Task>>('/tasks')
  const parsed = parseApiBody<Array<Task>>(response.data)
  return Array.isArray(parsed) ? parsed : []
}

export async function createTask(input: CreateTaskInput): Promise<Task> {
  const response = await api.post<Task>('/tasks', input)
  const parsed = parseApiBody<Task>(response.data)
  return parsed ?? response.data
}

export async function updateTask(
  id: string,
  input: UpdateTaskInput,
): Promise<Task> {
  const response = await api.put<Task>(`/tasks/${id}`, input)
  const parsed = parseApiBody<Task>(response.data)
  return parsed ?? response.data
}

export async function deleteTask(id: string): Promise<void> {
  await api.delete(`/tasks/${id}`)
}
