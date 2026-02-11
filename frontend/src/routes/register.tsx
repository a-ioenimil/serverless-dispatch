import { createFileRoute } from '@tanstack/react-router'

import SignUpPage from '../components/mvpblocks/signup-form-3'

export const Route = createFileRoute('/register')({
  component: SignUpPage,
})
