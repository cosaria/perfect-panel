import { createFileRoute } from '@tanstack/react-router';
import { parseUsersSearch } from '../../features/users/users-search';
import { UsersTableShell } from '../../features/users/users-table-shell';

export const Route = createFileRoute('/_authed/users')({
	component: UsersRoute,
	validateSearch: parseUsersSearch,
});

function UsersRoute() {
	const search = Route.useSearch();

	return <UsersTableShell search={search} />;
}
