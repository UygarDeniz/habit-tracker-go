import { createFileRoute, redirect } from '@tanstack/react-router'

import { z } from 'zod'
import GoogleAuthButton from '@/components/GoogleAuthButton'

const searchSchema = z.object({
  auth_error: z.string().optional(),
})

export const Route = createFileRoute('/auth/')({
  component: Auth,
  validateSearch: searchSchema,
  beforeLoad: ({ context }) => {
    if (context.user.isAuthenticated) {
      throw redirect({
        to: '/',
      })
    }
  },
})

function Auth() {
  const { auth_error } = Route.useSearch()

  return (
    <div className="flex flex-col items-center justify-center h-screen gap-4 text-center text-gray-500">
      <h1 className="text-2xl font-bold">StreakCraft</h1>
      <p className="text-sm">Login or sign up to continue</p>
      <div className="flex flex-col items-center justify-center gap-2">
        <GoogleAuthButton />
        {auth_error && <p className="text-red-500">{auth_error}</p>}
      </div>
    </div>
  )
}
