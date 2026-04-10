import { createFileRoute } from '@tanstack/react-router';
import { UsersTableShell } from '../../features/users/users-table-shell';

export const Route = createFileRoute('/_authed/users')({
	component: UsersRoute,
});

function UsersRoute() {
	return <UsersTableShell />;
}
