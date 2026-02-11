import { createFileRoute } from '@tanstack/react-router'

import SignInPage from '../components/mvpblocks/login-form-3'

export const Route = createFileRoute('/login')({
  component: SignInPage,
})
