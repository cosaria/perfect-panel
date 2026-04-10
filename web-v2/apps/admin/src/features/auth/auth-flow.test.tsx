import '@testing-library/jest-dom/vitest';

import { RouterContextProvider } from '@tanstack/react-router';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { createRouter } from '../../router';
import { parseResetPasswordSearch } from '../../routes/reset-password';
import {
	authStore,
	clearAuthSession,
	hydrateAuthSession,
	setAuthSession,
} from '../../shared/store/auth-store';
import { LoginForm } from './login-form';
import { ResetPasswordForm } from './reset-password-form';

describe('auth flow', () => {
	beforeEach(() => {
		window.sessionStorage.clear();
		clearAuthSession();
	});

	it('renders login fields and allows typing', () => {
		const router = createRouter();

		render(
			<RouterContextProvider router={router}>
				<LoginForm />
			</RouterContextProvider>,
		);

		expect(screen.getByText('邮箱')).toBeInTheDocument();
		expect(screen.getByText('密码')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument();

		fireEvent.change(screen.getByLabelText('邮箱'), { target: { value: 'admin@example.com' } });
		fireEvent.change(screen.getByLabelText('密码'), { target: { value: 'secret123' } });

		expect(screen.getByLabelText('邮箱')).toHaveValue('admin@example.com');
		expect(screen.getByLabelText('密码')).toHaveValue('secret123');
	});

	it('hydrates auth session from sessionStorage', () => {
		setAuthSession('token-123', 'admin@example.com');

		authStore.setState(() => ({
			accessToken: null,
			email: null,
			hydrated: false,
		}));

		hydrateAuthSession();

		expect(authStore.state.accessToken).toBe('token-123');
		expect(authStore.state.email).toBe('admin@example.com');
		expect(authStore.state.hydrated).toBe(true);
	});

	it('prefills reset token from route search', () => {
		const router = createRouter();
		const search = parseResetPasswordSearch({
			token: 'reset-token-123',
		});

		render(
			<RouterContextProvider router={router}>
				<ResetPasswordForm initialToken={search.token} />
			</RouterContextProvider>,
		);

		expect(screen.getByLabelText('令牌')).toHaveValue('reset-token-123');
	});
});
