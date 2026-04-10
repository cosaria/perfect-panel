import { createFileRoute } from '@tanstack/react-router';
import { LoginForm } from '../features/auth/login-form';
import { normalizeRedirectTarget } from '../shared/auth/redirect';

export function parseLoginSearch(search: Record<string, unknown>) {
	return {
		redirect: normalizeRedirectTarget(
			typeof search.redirect === 'string' ? search.redirect : undefined,
		),
	};
}

export const Route = createFileRoute('/login')({
	validateSearch: parseLoginSearch,
	component: LoginRoute,
});

function LoginRoute() {
	const { redirect } = Route.useSearch();

	return <LoginForm redirectTo={redirect} />;
}
