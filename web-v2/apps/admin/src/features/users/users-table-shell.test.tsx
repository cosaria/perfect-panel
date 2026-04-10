import '@testing-library/jest-dom/vitest';

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { clearAuthSession, setAuthSession } from '../../shared/store/auth-store';
import { useUsersListQuery } from './users-query';
import { UsersTableShell } from './users-table-shell';

const navigateMock = vi.fn();

vi.mock('@tanstack/react-router', async () => {
	const actual =
		await vi.importActual<typeof import('@tanstack/react-router')>('@tanstack/react-router');

	return {
		...actual,
		useNavigate: () => navigateMock,
	};
});

vi.mock('./users-query', () => ({
	useUsersListQuery: vi.fn(),
}));

describe('UsersTableShell', () => {
	beforeEach(() => {
		window.sessionStorage.clear();
		clearAuthSession();
		setAuthSession('token-123', 'admin@example.com');
		navigateMock.mockReset();
	});

	afterEach(() => {
		cleanup();
	});

	it('shows loading state while the list query is pending', () => {
		vi.mocked(useUsersListQuery).mockReturnValue({
			data: undefined,
			error: null,
			isError: false,
			isFetching: false,
			isLoading: true,
			isSuccess: false,
			refetch: vi.fn(),
		} as never);

		render(
			<UsersTableShell
				search={{
					order: 'asc',
					page: 1,
					pageSize: 10,
					q: '',
					sort: 'email',
				}}
			/>,
		);

		expect(screen.getByText('用户管理')).toBeInTheDocument();
		expect(document.querySelectorAll('.animate-pulse')).toHaveLength(5);
	});

	it('shows empty state when the query succeeds without rows', () => {
		vi.mocked(useUsersListQuery).mockReturnValue({
			data: {
				data: [],
				meta: {
					hasNext: false,
					hasPrev: false,
					page: 1,
					pageSize: 10,
					total: 0,
					totalPages: 0,
				},
			},
			error: null,
			isError: false,
			isFetching: false,
			isLoading: false,
			isSuccess: true,
			refetch: vi.fn(),
		} as never);

		render(
			<UsersTableShell
				search={{
					order: 'asc',
					page: 1,
					pageSize: 10,
					q: '',
					sort: 'email',
				}}
			/>,
		);

		expect(screen.getByText('没有找到匹配的用户')).toBeInTheDocument();
	});

	it('shows error state and refetch action when the query fails', async () => {
		const refetch = vi.fn();

		vi.mocked(useUsersListQuery).mockReturnValue({
			data: undefined,
			error: new Error('boom'),
			isError: true,
			isFetching: false,
			isLoading: false,
			isSuccess: false,
			refetch,
		} as never);

		render(
			<UsersTableShell
				search={{
					order: 'asc',
					page: 1,
					pageSize: 10,
					q: '',
					sort: 'email',
				}}
			/>,
		);

		fireEvent.click(screen.getByRole('button', { name: '重新加载' }));

		await waitFor(() => {
			expect(refetch).toHaveBeenCalled();
		});
	});

	it('updates route search when submitting a query, sorting, and paging', async () => {
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
					hasNext: true,
					hasPrev: false,
					page: 1,
					pageSize: 10,
					total: 11,
					totalPages: 2,
				},
			},
			error: null,
			isError: false,
			isFetching: false,
			isLoading: false,
			isSuccess: true,
			refetch: vi.fn(),
		} as never);

		render(
			<UsersTableShell
				search={{
					order: 'asc',
					page: 1,
					pageSize: 10,
					q: '',
					sort: 'email',
				}}
			/>,
		);

		fireEvent.change(screen.getByLabelText('搜索用户'), {
			target: { value: ' admin@example.com ' },
		});
		fireEvent.click(screen.getByRole('button', { name: '查询' }));

		await waitFor(() => {
			expect(navigateMock).toHaveBeenCalledWith({
				search: {
					order: 'asc',
					page: 1,
					pageSize: 10,
					q: 'admin@example.com',
					sort: 'email',
				},
				to: '/users',
			});
		});

		fireEvent.click(screen.getByRole('button', { name: /状态/ }));

		await waitFor(() => {
			expect(navigateMock).toHaveBeenCalledWith({
				search: {
					order: 'asc',
					page: 1,
					pageSize: 10,
					q: '',
					sort: 'status',
				},
				to: '/users',
			});
		});

		fireEvent.click(screen.getByRole('button', { name: '下一页' }));

		await waitFor(() => {
			expect(navigateMock).toHaveBeenCalledWith({
				search: {
					order: 'asc',
					page: 2,
					pageSize: 10,
					q: '',
					sort: 'email',
				},
				to: '/users',
			});
		});
	});
});
