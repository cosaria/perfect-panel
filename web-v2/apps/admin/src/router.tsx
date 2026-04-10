import { QueryClient } from '@tanstack/react-query'
import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { Route as indexRoute } from './routes/index'
import { Route as rootRoute } from './routes/__root'

export interface AppRouterContext {
	queryClient: QueryClient
}

export const routeTree = rootRoute.addChildren([indexRoute])

export function createRouter() {
	const queryClient = new QueryClient()

	return createTanStackRouter({
		context: {
			queryClient,
		},
		defaultPreload: 'intent',
		routeTree,
		scrollRestoration: true,
	})
}

declare module '@tanstack/react-router' {
	interface Register {
		router: ReturnType<typeof createRouter>
	}
}
