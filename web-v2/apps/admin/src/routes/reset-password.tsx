import { createFileRoute } from '@tanstack/react-router';
import { ResetPasswordForm } from '../features/auth/reset-password-form';

export function parseResetPasswordSearch(search: Record<string, unknown>) {
	return {
		token: typeof search.token === 'string' ? search.token : '',
	};
}

export const Route = createFileRoute('/reset-password')({
	validateSearch: parseResetPasswordSearch,
	component: ResetPasswordRoute,
});

function ResetPasswordRoute() {
	const { token } = Route.useSearch();

	return <ResetPasswordForm initialToken={token} />;
}
