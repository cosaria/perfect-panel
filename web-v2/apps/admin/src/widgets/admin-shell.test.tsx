import '@testing-library/jest-dom/vitest';

import { RouterContextProvider, createMemoryHistory } from '@tanstack/react-router';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useUsersListQuery } from '../features/users/users-query';
import { UsersTableShell } from '../features/users/users-table-shell';
import { createRouter } from '../router';
import { clearAuthSession, setAuthSession } from '../shared/store/auth-store';
import { AdminSidebar } from './admin-sidebar';
import { AdminTopbar } from './admin-topbar';

vi.mock('../features/users/users-query', () => ({
	useUsersListQuery: vi.fn(),
}));

describe('AdminSidebar', () => {
	beforeEach(() => {
		window.sessionStorage.clear();
		clearAuthSession();
		window.scrollTo = vi.fn();
	});

	afterEach(() => {
		cleanup();
	});

	it('renders the first-stage navigation items', () => {
		const router = createRouter();

		render(
			<RouterContextProvider router={router}>
				<AdminSidebar />
			</RouterContextProvider>,
		);

		expect(screen.getByText('Dashboard')).toBeInTheDocument();
		expect(screen.getByText('用户管理')).toBeInTheDocument();
		expect(screen.getByText('系统设置')).toBeInTheDocument();
	});

	it('renders the authed users route with active navigation and topbar identity', async () => {
		setAuthSession('token-123', 'admin@example.com');
		vi.mocked(useUsersListQuery).mockReturnValue({
			data: {
				data: [
					{
						email: 'admin@example.com',
						id: '00000000-0000-0000-0000-000000000001',
						status: 'active',
					},
				],
				meta: {
					hasNext: false,
					hasPrev: false,
					page: 1,
					pageSize: 10,
					total: 1,
					totalPages: 1,
				},
			},
			error: null,
			isError: false,
			isFetching: false,
			isLoading: false,
			isSuccess: true,
			refetch: vi.fn(),
		} as never);

		const history = createMemoryHistory({
			initialEntries: ['/users'],
		});
		const router = createRouter({ history });

		render(
			<RouterContextProvider router={router}>
				<div className="min-h-screen bg-slate-950 text-slate-50 lg:grid lg:grid-cols-[18rem_minmax(0,1fr)]">
					<AdminSidebar />
					<div className="min-w-0">
						<AdminTopbar />
						<main className="p-6 lg:p-8">
							<UsersTableShell
								search={{
									order: 'asc',
									page: 1,
									pageSize: 10,
									q: '',
									sort: 'email',
								}}
							/>
						</main>
					</div>
				</div>
			</RouterContextProvider>,
		);

		await waitFor(() => {
			expect(screen.getByText('真实用户列表')).toBeInTheDocument();
		});

		expect(screen.getAllByText('admin@example.com')).toHaveLength(2);
		expect(
			screen
				.getAllByRole('link', { name: /用户管理/ })
				.find((link) => link.getAttribute('aria-current') === 'page'),
		).toHaveClass('border-cyan-300/40');
		expect(screen.getByRole('link', { name: /Dashboard/ })).not.toHaveClass('border-cyan-300/40');
	});
});
