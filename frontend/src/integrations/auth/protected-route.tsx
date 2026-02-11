import { ReactNode } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useAuth } from './auth-context'

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#0b0b0d]">
        <div className="text-center">
          <div className="mb-4 h-8 w-8 animate-spin rounded-full border-4 border-amber-200 border-t-transparent" />
          <p className="text-sm text-amber-200/60">Loading...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    navigate({ to: '/login', replace: true })
    return null
  }

  return <>{children}</>
}
