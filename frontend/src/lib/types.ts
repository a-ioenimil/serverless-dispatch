export type TaskStatus = 'OPEN' | 'IN_PROGRESS' | 'DONE'

export type TaskPriority = 'LOW' | 'MEDIUM' | 'HIGH'

export interface Task {
  id: string
  title: string
  status: TaskStatus
  priority: TaskPriority
  assignee_id?: string | null
  created_by: string
  created_at: string
}

export interface CreateTaskInput {
  title: string
  description?: string
  priority: TaskPriority
  assignee_id?: string | null
}

export interface UpdateTaskInput {
  status?: TaskStatus
  priority?: TaskPriority
  assignee_id?: string | null
}

export interface UserSummary {
  id: string
  username: string
  email: string
  role?: string
  status?: string
}
