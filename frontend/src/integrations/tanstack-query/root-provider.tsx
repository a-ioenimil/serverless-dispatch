import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

export function getContext() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 15000,
        refetchOnWindowFocus: false,
        retry: (failureCount, error) => {
          const status = (error as { response?: { status?: number } }).response
            ?.status
          if (!status) {
            return failureCount < 2
          }
          return status >= 500 && failureCount < 2
        },
        retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 3000),
      },
      mutations: {
        retry: (failureCount, error) => {
          const status = (error as { response?: { status?: number } }).response
            ?.status
          if (!status) {
            return failureCount < 1
          }
          return status >= 500 && failureCount < 1
        },
        retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 2000),
      },
    },
  })
  return {
    queryClient,
  }
}

export function Provider({
  children,
  queryClient,
}: {
  children: React.ReactNode
  queryClient: QueryClient
}) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
