import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from 'react'
import { getStoredTokens, getStoredUser, signOut as signOutAuth, isTokenExpired } from '../../lib/auth'
import type { AuthTokens, AuthUser } from '../../lib/auth'

interface AuthContextType {
  user: AuthUser | null
  tokens: AuthTokens | null
  isLoading: boolean
  isAuthenticated: boolean
  signOut: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [tokens, setTokens] = useState<AuthTokens | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Initialize auth state from localStorage
  useEffect(() => {
    const storedTokens = getStoredTokens()
    const storedUser = getStoredUser()

    if (storedTokens && storedUser) {
      // Check if token is expired
      if (isTokenExpired(storedTokens.accessToken)) {
        // Token expired, clear storage
        signOutAuth()
        setUser(null)
        setTokens(null)
      } else {
        setTokens(storedTokens)
        setUser(storedUser)
      }
    }

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
