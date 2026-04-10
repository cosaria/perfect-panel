import { createRoute } from '@tanstack/react-router'
import { Route as rootRoute } from './__root'

export const Route = createRoute({
	getParentRoute: () => rootRoute,
	path: '/',
	component: AdminShellRoute,
})

function AdminShellRoute() {
	return (
		<main>
			<h1>admin-shell</h1>
		</main>
	)
}
