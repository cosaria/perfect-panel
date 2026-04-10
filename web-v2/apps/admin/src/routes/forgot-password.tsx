import { createFileRoute } from '@tanstack/react-router';
import { ForgotPasswordForm } from '../features/auth/forgot-password-form';

export const Route = createFileRoute('/forgot-password')({
	component: ForgotPasswordRoute,
});

function ForgotPasswordRoute() {
	return <ForgotPasswordForm />;
}
