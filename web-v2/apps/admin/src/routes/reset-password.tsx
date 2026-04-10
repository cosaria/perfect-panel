import { createFileRoute } from '@tanstack/react-router';
import { ResetPasswordForm } from '../features/auth/reset-password-form';

export const Route = createFileRoute('/reset-password')({
	component: ResetPasswordRoute,
});

function ResetPasswordRoute() {
	return <ResetPasswordForm />;
}
