import { QueryClient } from '@tanstack/react-query'
import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'

export interface AppRouterContext {
	queryClient: QueryClient
}

export function getRouter() {
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

export function createRouter() {
	return getRouter()
}

declare module '@tanstack/react-router' {
	interface Register {
		router: ReturnType<typeof getRouter>
	}
}
