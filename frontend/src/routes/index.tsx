import { createFileRoute } from '@tanstack/react-router'

import { TaskBoard } from '../components/kanban/TaskBoard'

export const Route = createFileRoute('/')({
  component: TaskBoard,
})
