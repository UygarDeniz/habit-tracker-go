import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import TanStackQueryLayout from '../integrations/tanstack-query/layout.tsx'
import Header from '../components/Header'

import type { QueryClient } from '@tanstack/react-query'

import type { UserContextType } from '@/contexts/UserContext.tsx'

interface MyRouterContext {
  queryClient: QueryClient
  user: UserContextType
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: () => (
    <>
      <Header />
      <Outlet />
      <TanStackRouterDevtools />

      <TanStackQueryLayout />
    </>
  ),
})
