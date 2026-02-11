'use client'

import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { Loader2, ShieldCheck, Sparkles } from 'lucide-react'

import { Button } from '../ui/button'
import { InputOTP, InputOTPGroup, InputOTPSlot } from '../ui/input-otp'
import { confirmSignUp, resendSignUpCode } from '../../lib/auth'

const OTP_LENGTH = 6

export default function VerifyOtpPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [otp, setOtp] = useState('')
  const [loading, setLoading] = useState(false)
  const [resending, setResending] = useState(false)

  const maskedEmail = useMemo(() => {
    if (!email) return ''
    const [localPart, domain] = email.split('@')
    if (!domain) return email
    const visible = localPart.slice(0, 2)
    return `${visible}${'*'.repeat(Math.max(localPart.length - 2, 0))}@${domain}`
  }, [email])

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const queryEmail = params.get('email')?.trim() || ''
    const storedEmail = localStorage.getItem('pending_signup_email') || ''
    setEmail(queryEmail || storedEmail)
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!email.trim()) {
      toast.error('Missing email. Please sign up again.')
      return
    }

    if (otp.length !== OTP_LENGTH) {
      toast.error('Enter the 6-digit verification code')
      return
    }

    setLoading(true)

    try {
      await confirmSignUp(email.trim(), otp)
      toast.success('Verification complete. You can now sign in.')
      localStorage.removeItem('pending_signup_email')
      navigate({ to: '/login', replace: true })
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Verification failed'
      toast.error(message)
      console.error('Verification error:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleResend = async () => {
    if (!email.trim()) {
      toast.error('Missing email. Please sign up again.')
      return
    }

    setResending(true)

    try {
      const result = await resendSignUpCode(email.trim())
      const destination = result.destination ? ` to ${result.destination}` : ''
      toast.success(`Verification code resent${destination}.`)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Resend failed'
      toast.error(message)
      console.error('Resend error:', error)
    } finally {
      setResending(false)
    }
  }

  return (
    <div className="relative flex min-h-screen w-full items-center justify-center overflow-hidden bg-[#0b0b0d] p-4">
      <style>{`
        :root {
          --auth-amber: #fcd34d;
          --auth-ember: #f59e0b;
          --auth-ink: #111114;
        }
        .auth-grid {
          background-image:
            linear-gradient(120deg, rgba(252, 211, 77, 0.08), transparent 55%),
            radial-gradient(circle at 20% 20%, rgba(252, 211, 77, 0.15), transparent 60%),
            radial-gradient(circle at 80% 80%, rgba(245, 158, 11, 0.2), transparent 55%);
        }
        .auth-hero {
          background-image:
            linear-gradient(140deg, rgba(10, 10, 12, 0.7), rgba(10, 10, 12, 0.2)),
            url('/assets/images/login_side_backgrond_image.jpg');
          background-size: cover;
          background-position: center;
        }
        .glow-card {
          box-shadow:
            0 40px 120px rgba(0, 0, 0, 0.55),
            inset 0 0 0 1px rgba(255, 255, 255, 0.06);
        }
        .pulse-dot {
          animation: pulse 2.4s ease-in-out infinite;
        }
        @keyframes pulse {
          0%, 100% { transform: scale(1); opacity: 0.7; }
          50% { transform: scale(1.4); opacity: 1; }
        }
        .fade-rise {
          animation: fadeRise 0.7s ease forwards;
        }
        @keyframes fadeRise {
          from { opacity: 0; transform: translateY(10px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .shine-btn {
          position: relative;
          overflow: hidden;
        }
        .shine-btn::after {
          content: '';
          position: absolute;
          top: 0;
          left: -120%;
          width: 60%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255,255,255,0.35), transparent);
          transform: skewX(-12deg);
          transition: left 0.6s ease;
        }
        .shine-btn:hover::after { left: 140%; }
      `}</style>

      <div className="absolute inset-0 auth-grid opacity-90" />

      <div className="relative z-10 w-full max-w-5xl">
        <div className="glow-card overflow-hidden rounded-[32px] border border-white/10 bg-[#111114]/80">
          <div className="grid min-h-[620px] lg:grid-cols-[1.05fr_0.95fr]">
            <div className="auth-hero relative flex flex-col justify-between gap-10 p-10 text-white">
              <div className="space-y-6">
                <div className="inline-flex items-center gap-3 rounded-full border border-white/15 bg-white/5 px-4 py-2 text-xs uppercase tracking-[0.3em] text-amber-200/80">
                  <span className="pulse-dot h-2 w-2 rounded-full bg-amber-300" />
                  Verify Secure Access
                </div>
                <div className="space-y-4">
                  <h1 className="text-4xl font-semibold leading-tight">
                    Dispatch Console
                    <span className="text-amber-300">.</span>
                  </h1>
                  <p className="max-w-md text-sm text-slate-200/70">
                    Confirm your email with the 6-digit code we sent to complete
                    onboarding.
                  </p>
                </div>
              </div>

              <div className="space-y-4 text-sm text-slate-200/70">
                <div className="flex items-center gap-3">
                  <ShieldCheck className="h-5 w-5 text-amber-200" />
                  One-time code expires quickly
                </div>
                <div className="flex items-center gap-3">
                  <Sparkles className="h-5 w-5 text-amber-200" />
                  Complete setup to unlock your dashboard
                </div>
              </div>
            </div>

            <div className="flex flex-col justify-center border-t border-white/5 bg-black/30 p-10 lg:border-t-0 lg:border-l">
              <div className="mx-auto w-full max-w-md space-y-8">
                <div className="fade-rise space-y-2 text-center">
                  <h2 className="text-3xl font-semibold text-amber-50">
                    Verify your email
                  </h2>
                  <p className="text-sm text-slate-200/60">
                    Enter the 6-digit code sent to {maskedEmail || 'your inbox'}
                    .
                  </p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-6">
                  <div className="space-y-3">
                    <label className="text-xs uppercase tracking-[0.2em] text-amber-100/70">
                      Verification code
                    </label>
                    <InputOTP
                      maxLength={OTP_LENGTH}
                      value={otp}
                      onChange={setOtp}
                      inputMode="numeric"
                      pattern="\d*"
                      className="justify-center"
                      containerClassName="justify-center"
                    >
                      <InputOTPGroup>
                        {Array.from({ length: OTP_LENGTH }).map((_, index) => (
                          <InputOTPSlot
                            key={index}
                            index={index}
                            className="h-12 w-12 border-white/10 bg-black/40 text-amber-50 text-base"
                          />
                        ))}
                      </InputOTPGroup>
                    </InputOTP>
                  </div>

                  <Button
                    type="submit"
                    className="shine-btn h-11 w-full bg-amber-200 text-black shadow-[0_12px_40px_rgba(252,211,77,0.35)] hover:bg-amber-200"
                    disabled={loading}
                  >
                    {loading ? (
                      <span className="flex items-center gap-2">
                        <Loader2 className="h-4 w-4 animate-spin" />
                        Verifying...
                      </span>
                    ) : (
                      'Verify code'
                    )}
                  </Button>
                </form>

                <div className="flex flex-col items-center gap-3 text-sm text-amber-100/70">
                  <button
                    type="button"
                    onClick={handleResend}
                    disabled={resending}
                    className="text-amber-200/80 transition hover:text-amber-200 disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    {resending ? 'Resending...' : 'Resend code'}
                  </button>
                  <div>
                    Back to{' '}
                    <Link
                      to="/login"
                      className="text-amber-200/80 transition hover:text-amber-200"
                    >
                      sign in
                    </Link>
                  </div>
                </div>

                <div className="text-center text-xs text-amber-100/50">
                  Use the latest code sent to your email address.
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
