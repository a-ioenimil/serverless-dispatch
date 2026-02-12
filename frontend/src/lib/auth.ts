import {
  CognitoIdentityProviderClient,
  ConfirmSignUpCommand,
  InitiateAuthCommand,
  ResendConfirmationCodeCommand,
  SignUpCommand,
} from '@aws-sdk/client-cognito-identity-provider'

const clientId = import.meta.env.VITE_COGNITO_CLIENT_ID as string
const userPoolId = import.meta.env.VITE_COGNITO_USER_POOL_ID as string
const region = (import.meta.env.VITE_AWS_REGION as string) || 'us-east-1'

const client = new CognitoIdentityProviderClient({ region })

export interface AuthTokens {
  idToken: string
  accessToken: string
  refreshToken?: string
  expiresIn: number
}

export interface AuthUser {
  id: string
  email: string
  groups: string[]
}

export interface SignUpResult {
  userSub: string
  userConfirmed: boolean
}

export interface ConfirmSignUpResult {
  userConfirmed: boolean
}

export interface ResendConfirmationResult {
  destination?: string
}

const TOKEN_STORAGE_KEY = 'auth_tokens'
const USER_STORAGE_KEY = 'auth_user'

/**
 * Sign in with email and password
 */
export async function signIn(
  email: string,
  password: string,
): Promise<AuthTokens> {
  try {
    const command = new InitiateAuthCommand({
      ClientId: clientId,
      AuthFlow: 'USER_PASSWORD_AUTH',
      AuthParameters: {
        USERNAME: email,
        PASSWORD: password,
      },
    })

    const response = await client.send(command)

    if (!response.AuthenticationResult) {
      throw new Error('No authentication result returned')
    }

    const tokens: AuthTokens = {
      idToken: response.AuthenticationResult.IdToken || '',
      accessToken: response.AuthenticationResult.AccessToken || '',
      refreshToken: response.AuthenticationResult.RefreshToken,
      expiresIn: response.AuthenticationResult.ExpiresIn || 3600,
    }

    // Store tokens
    localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(tokens))

    // Extract user info from id token
    const user = extractUserFromToken(tokens.idToken)
    if (user) {
      localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(user))
    }

    return tokens
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Sign in failed'
    throw new Error(message)
  }
}

/**
 * Sign up with email and password
 */
export async function signUp(
  email: string,
  password: string,
): Promise<SignUpResult> {
  try {
    const command = new SignUpCommand({
      ClientId: clientId,
      Username: email,
      Password: password,
      UserAttributes: [
        {
          Name: 'email',
          Value: email,
        },
      ],
    })

    const response = await client.send(command)

    return {
      userSub: response.UserSub || '',
      userConfirmed: Boolean(response.UserConfirmed),
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Sign up failed'
    throw new Error(message)
  }
}

/**
 * Confirm sign up with the verification code
 */
export async function confirmSignUp(
  email: string,
  code: string,
): Promise<ConfirmSignUpResult> {
  try {
    const command = new ConfirmSignUpCommand({
      ClientId: clientId,
      Username: email,
      ConfirmationCode: code,
    })

    await client.send(command)

    return {
      userConfirmed: true,
    }
  } catch (error) {
    const message =
      error instanceof Error ? error.message : 'Verification failed'
    throw new Error(message)
  }
}

/**
 * Resend verification code for sign up
 */
export async function resendSignUpCode(
  email: string,
): Promise<ResendConfirmationResult> {
  try {
    const command = new ResendConfirmationCodeCommand({
      ClientId: clientId,
      Username: email,
    })

    const response = await client.send(command)

    return {
      destination: response.CodeDeliveryDetails?.Destination,
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Resend failed'
    throw new Error(message)
  }
}

function decodeJwtPayload(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    const payload = parts[1]
      .replace(/-/g, '+')
      .replace(/_/g, '/')
      .padEnd(Math.ceil(parts[1].length / 4) * 4, '=')

    return JSON.parse(atob(payload)) as Record<string, unknown>
  } catch {
    return null
  }
}

/**
 * Extract user info from ID token
 */
export function extractUserFromToken(token: string): AuthUser | null {
  const decoded = decodeJwtPayload(token)
  if (!decoded) return null

  const groupsValue = decoded['cognito:groups']
  const groups = Array.isArray(groupsValue)
    ? groupsValue.filter(Boolean).map(String)
    : typeof groupsValue === 'string'
      ? groupsValue.split(',')
      : []

  return {
    id: String(decoded.sub || ''),
    email: String(decoded.email || ''),
    groups,
  }
}

export function deriveUserFromTokens(
  tokens: AuthTokens | null,
): AuthUser | null {
  if (!tokens?.idToken) return null
  return extractUserFromToken(tokens.idToken)
}

/**
 * Get stored tokens from localStorage
 */
export function getStoredTokens(): AuthTokens | null {
  const stored = localStorage.getItem(TOKEN_STORAGE_KEY)
  return stored ? JSON.parse(stored) : null
}

/**
 * Get stored user from localStorage
 */
export function getStoredUser(): AuthUser | null {
  const stored = localStorage.getItem(USER_STORAGE_KEY)
  if (stored) return JSON.parse(stored)

  const tokens = getStoredTokens()
  return deriveUserFromTokens(tokens)
}

/**
 * Sign out and clear tokens
 */
export function signOut(): void {
  localStorage.removeItem(TOKEN_STORAGE_KEY)
  localStorage.removeItem(USER_STORAGE_KEY)
}

/**
 * Check if user is authenticated
 */
export function isAuthenticated(): boolean {
  const tokens = getStoredTokens()
  return !!tokens?.accessToken
}

/**
 * Check if token is expired
 */
export function isTokenExpired(token: string): boolean {
  const decoded = decodeJwtPayload(token)
  if (!decoded || typeof decoded.exp !== 'number') return true

  const expiresAt = decoded.exp * 1000 // Convert to milliseconds
  return Date.now() >= expiresAt
}
