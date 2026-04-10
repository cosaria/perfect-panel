import { createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '../shared/auth/require-auth';

export const Route = createFileRoute('/')({
	component: AdminShellRoute,
});

function AdminShellRoute() {
	return (
		<RequireAuth>
			<main>
				<h1>admin-shell</h1>
			</main>
		</RequireAuth>
	);
}
