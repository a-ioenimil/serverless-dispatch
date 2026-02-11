import { createFileRoute } from '@tanstack/react-router'

import { TaskBoard } from '../components/kanban/TaskBoard'
import { ProtectedRoute } from '../integrations/auth/protected-route'

export const Route = createFileRoute('/')({
  component: () => (
    <ProtectedRoute>
      <TaskBoard />
    </ProtectedRoute>
  ),
})
