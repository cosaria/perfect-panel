import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
	component: AdminShellRoute,
});

function AdminShellRoute() {
	return (
		<main>
			<h1>admin-shell</h1>
		</main>
	);
}
