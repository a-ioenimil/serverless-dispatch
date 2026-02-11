import { createFileRoute } from '@tanstack/react-router'

import VerifyOtpPage from '../components/mvpblocks/verify-otp-form-3'

export const Route = createFileRoute('/verify-otp')({
  component: VerifyOtpPage,
})
