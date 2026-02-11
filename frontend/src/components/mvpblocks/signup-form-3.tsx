'use client'

import { useState } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import {
  Eye,
  EyeOff,
  Loader2,
  Lock,
  Mail,
  ShieldCheck,
  Sparkles,
} from 'lucide-react'

import { Button } from '../ui/button'
import { Checkbox } from '../ui/checkbox'
import { Input } from '../ui/input'
import { signUp } from '../../lib/auth'

export default function SignUpPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [acceptTerms, setAcceptTerms] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!acceptTerms) {
      toast.error('Please accept the terms to continue')
      return
    }

    setLoading(true)

    try {
      await signUp(email, password)
      toast.success('Account created. Check your email to verify.')
      localStorage.setItem('pending_signup_email', email)
      const encodedEmail = encodeURIComponent(email)
      navigate({ to: `/verify-otp?email=${encodedEmail}`, replace: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Sign up failed'
      toast.error(message)
      console.error('Sign up error:', error)
    } finally {
      setLoading(false)
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
                  Cognito Secure Access
                </div>
                <div className="space-y-4">
                  <h1 className="text-4xl font-semibold leading-tight">
                    Dispatch Console
                    <span className="text-amber-300">.</span>
                  </h1>
                  <p className="max-w-md text-sm text-slate-200/70">
                    Create your account to orchestrate tasks, monitor workflow
                    health, and stay synced with your serverless stack.
                  </p>
                </div>
              </div>

              <div className="space-y-4 text-sm text-slate-200/70">
                <div className="flex items-center gap-3">
                  <ShieldCheck className="h-5 w-5 text-amber-200" />
                  Email verification required for access
                </div>
                <div className="flex items-center gap-3">
                  <Sparkles className="h-5 w-5 text-amber-200" />
                  Secure onboarding powered by Cognito
                </div>
              </div>
            </div>

            <div className="flex flex-col justify-center border-t border-white/5 bg-black/30 p-10 lg:border-t-0 lg:border-l">
              <div className="mx-auto w-full max-w-md space-y-8">
                <div className="fade-rise space-y-2 text-center">
                  <h2 className="text-3xl font-semibold text-amber-50">
                    Create your account
                  </h2>
                  <p className="text-sm text-slate-200/60">
                    Use your email and a secure password to get started.
                  </p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-6">
                  <div className="space-y-2">
                    <label
                      htmlFor="email"
                      className="text-xs uppercase tracking-[0.2em] text-amber-100/70"
                    >
                      Email address
                    </label>
                    <div className="relative">
                      <Mail className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-amber-100/40" />
                      <Input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        placeholder="you@company.com"
                        className="h-11 border-white/10 bg-black/40 pl-10 text-amber-50 placeholder:text-amber-100/30 focus-visible:ring-amber-200/30"
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <label
                      htmlFor="password"
                      className="text-xs uppercase tracking-[0.2em] text-amber-100/70"
                    >
                      Password
                    </label>
                    <div className="relative">
                      <Lock className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-amber-100/40" />
                      <Input
                        id="password"
                        type={showPassword ? 'text' : 'password'}
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        placeholder="Create a strong password"
                        className="h-11 border-white/10 bg-black/40 pl-10 pr-12 text-amber-50 placeholder:text-amber-100/30 focus-visible:ring-amber-200/30"
                      />
                      <button
                        type="button"
                        aria-label={
                          showPassword ? 'Hide password' : 'Show password'
                        }
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-amber-100/60 transition hover:text-amber-100"
                        onClick={() => setShowPassword(!showPassword)}
                      >
                        {showPassword ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </button>
                    </div>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <label className="flex items-center gap-2 text-amber-100/70">
                      <Checkbox
                        checked={acceptTerms}
                        onCheckedChange={(checked) =>
                          setAcceptTerms(checked === true)
                        }
                        className="border-white/20 data-[state=checked]:bg-amber-200 data-[state=checked]:text-black"
                      />
                      I agree to the terms
                    </label>
                    <Link
                      to="/login"
                      className="text-amber-200/80 transition hover:text-amber-200"
                    >
                      Already have an account?
                    </Link>
                  </div>

                  <Button
                    type="submit"
                    className="shine-btn h-11 w-full bg-amber-200 text-black shadow-[0_12px_40px_rgba(252,211,77,0.35)] hover:bg-amber-200"
                    disabled={loading}
                  >
                    {loading ? (
                      <span className="flex items-center gap-2">
                        <Loader2 className="h-4 w-4 animate-spin" />
                        Creating account...
                      </span>
                    ) : (
                      'Sign up'
                    )}
                  </Button>
                </form>

                <div className="text-center text-xs text-amber-100/50">
                  Verify your email to activate your Dispatch Console account.
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
