import '@testing-library/jest-dom/vitest';

import { RouterContextProvider, createMemoryHistory } from '@tanstack/react-router';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { publicAuthSessionCreate } from '@web-v2/api-client';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createRouter } from '../../router';
import { parseLoginSearch } from '../../routes/login';
import { parseResetPasswordSearch } from '../../routes/reset-password';
import { AuthHydrationFallback, RequireAuth } from '../../shared/auth/require-auth';
import {
	authStore,
	clearAuthSession,
	hydrateAuthSession,
	setAuthSession,
} from '../../shared/store/auth-store';
import { LoginForm } from './login-form';
import { ResetPasswordForm } from './reset-password-form';

const navigateMock = vi.fn();
const navigateComponentMock = vi.fn();
const routerLocationMock = {
	hash: '',
	href: '/',
	pathname: '/',
	searchStr: '',
};

vi.mock('@tanstack/react-router', async () => {
	const actual =
		await vi.importActual<typeof import('@tanstack/react-router')>('@tanstack/react-router');

	return {
		...actual,
		Navigate: (props: unknown) => {
			navigateComponentMock(props);

			return null;
		},
		useNavigate: () => navigateMock,
		useRouterState: (options?: {
			select?: (state: { location: typeof routerLocationMock }) => unknown;
		}) =>
			options?.select
				? options.select({ location: routerLocationMock })
				: { location: routerLocationMock },
	};
});

vi.mock('@web-v2/api-client', () => ({
	publicAuthSessionCreate: vi.fn(),
}));

describe('auth flow', () => {
	beforeEach(() => {
		window.sessionStorage.clear();
		clearAuthSession();
		window.scrollTo = vi.fn();
		navigateComponentMock.mockReset();
		navigateMock.mockReset();
		navigateMock.mockResolvedValue(undefined);
		routerLocationMock.hash = '';
		routerLocationMock.href = '/';
		routerLocationMock.pathname = '/';
		routerLocationMock.searchStr = '';
		vi.mocked(publicAuthSessionCreate).mockReset();
	});

	afterEach(() => {
		cleanup();
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

	it('accepts only safe internal login redirect targets', () => {
		expect(parseLoginSearch({ redirect: '/users?status=active' })).toEqual({
			redirect: '/users?status=active',
		});
		expect(parseLoginSearch({ redirect: 'https://example.com/admin' })).toEqual({
			redirect: undefined,
		});
	});

	it('shows a hydration placeholder before the session restore completes', () => {
		render(<AuthHydrationFallback />);

		expect(screen.getByText('正在恢复登录状态…')).toBeInTheDocument();
	});

	it('preserves protected-route redirect target for unauthenticated access', async () => {
		authStore.setState(() => ({
			accessToken: null,
			email: null,
			hydrated: false,
		}));
		routerLocationMock.href = '/users';
		routerLocationMock.pathname = '/users';

		render(
			<RequireAuth>
				<div>secret</div>
			</RequireAuth>,
		);

		await waitFor(() => {
			expect(navigateComponentMock).toHaveBeenCalledWith({
				replace: true,
				search: {
					redirect: '/users',
				},
				to: '/login',
			});
		});
	});

	it('redirects to the preserved target after login succeeds', async () => {
		vi.mocked(publicAuthSessionCreate).mockResolvedValue({
			data: {
				data: {
					accessToken: 'token-123',
				},
			},
		} as Awaited<ReturnType<typeof publicAuthSessionCreate>>);

		const history = createMemoryHistory({
			initialEntries: ['/login'],
		});
		const router = createRouter({ history });

		render(
			<RouterContextProvider router={router}>
				<LoginForm redirectTo="/users" />
			</RouterContextProvider>,
		);

		fireEvent.change(screen.getByLabelText('邮箱'), {
			target: { value: 'admin@example.com' },
		});
		fireEvent.change(screen.getByLabelText('密码'), { target: { value: 'secret123' } });
		fireEvent.click(screen.getByRole('button', { name: '登录' }));

		await waitFor(() => {
			expect(navigateMock).toHaveBeenCalledWith({
				href: '/users',
				replace: true,
			});
		});
	});
});
