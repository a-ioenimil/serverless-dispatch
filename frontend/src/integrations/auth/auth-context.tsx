import { createContext, useContext, useEffect, useState } from 'react'
import {
  deriveUserFromTokens,
  getStoredTokens,
  getStoredUser,
  isTokenExpired,
  signOut as signOutAuth,
} from '../../lib/auth'
import type { ReactNode } from 'react'
import type { AuthTokens, AuthUser } from '../../lib/auth'

interface AuthContextType {
  user: AuthUser | null
  tokens: AuthTokens | null
  isLoading: boolean
  isAuthenticated: boolean
  signOut: () => void
  syncFromStorage: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [tokens, setTokens] = useState<AuthTokens | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const syncFromStorage = () => {
    const storedTokens = getStoredTokens()
    const storedUser = getStoredUser()

    if (storedTokens && storedUser) {
      if (isTokenExpired(storedTokens.accessToken)) {
        signOutAuth()
        setUser(null)
        setTokens(null)
        return
      }

      setTokens(storedTokens)
      setUser(storedUser)
      return
    }

    if (storedTokens) {
      if (isTokenExpired(storedTokens.accessToken)) {
        signOutAuth()
        setUser(null)
        setTokens(null)
        return
      }

      const derivedUser = deriveUserFromTokens(storedTokens)
      setTokens(storedTokens)
      setUser(derivedUser)
      return
    }

    setUser(null)
    setTokens(null)
  }

  // Initialize auth state from localStorage
  useEffect(() => {
    syncFromStorage()
    setIsLoading(false)
  }, [])

  const handleSignOut = () => {
    signOutAuth()
    setUser(null)
    setTokens(null)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        tokens,
        isLoading,
        isAuthenticated: !!tokens?.accessToken && !!user,
        signOut: handleSignOut,
        syncFromStorage,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
