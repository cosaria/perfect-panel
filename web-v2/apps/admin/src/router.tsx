import { QueryClient } from '@tanstack/react-query';
import { type RouterHistory, createRouter as createTanStackRouter } from '@tanstack/react-router';
import { routeTree } from './routeTree.gen';

export interface AppRouterContext {
	queryClient: QueryClient;
}

export function getRouter(history?: RouterHistory) {
	const queryClient = new QueryClient();

	return createTanStackRouter({
		context: {
			queryClient,
		},
		defaultPreload: 'intent',
		history,
		routeTree,
		scrollRestoration: true,
	});
}

export function createRouter(options?: { history?: RouterHistory }) {
	return getRouter(options?.history);
}

declare module '@tanstack/react-router' {
	interface Register {
		router: ReturnType<typeof getRouter>;
	}
}
