import {
	HeadContent,
	Outlet,
	Scripts,
	createRootRouteWithContext,
} from '@tanstack/react-router'
import { AppProviders } from '../app/providers'
import type { AppRouterContext } from '../router'

export const Route = createRootRouteWithContext<AppRouterContext>()({
	component: RootComponent,
})

function RootComponent() {
	return (
		<html lang="zh-CN">
			<head>
				<meta charSet="utf-8" />
				<meta
					name="viewport"
					content="width=device-width, initial-scale=1, viewport-fit=cover"
				/>
				<HeadContent />
			</head>
			<body>
				<AppProviders>
					<Outlet />
				</AppProviders>
				<Scripts />
			</body>
		</html>
	)
}
