import { adminAccessUserList } from '@web-v2/api-client';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createUsersListQueryKey, fetchUsersList } from './users-query';

vi.mock('@web-v2/api-client', () => ({
	adminAccessUserList: vi.fn(),
}));

describe('users query', () => {
	beforeEach(() => {
		vi.mocked(adminAccessUserList).mockReset();
	});

	it('creates a stable query key from route search state', () => {
		expect(
			createUsersListQueryKey({
				order: 'asc',
				page: 2,
				pageSize: 20,
				q: 'admin',
				sort: 'status',
			}),
		).toEqual([
			'admin',
			'users',
			{
				order: 'asc',
				page: 2,
				pageSize: 20,
				q: 'admin',
				sort: 'status',
			},
		]);
	});

	it('sends search params and authorization header to the SDK', async () => {
		vi.mocked(adminAccessUserList).mockResolvedValue({
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
		} as unknown as Awaited<ReturnType<typeof adminAccessUserList>>);

		await fetchUsersList(
			{
				order: 'desc',
				page: 3,
				pageSize: 20,
				q: 'admin',
				sort: 'email',
			},
			'token-123',
		);

		expect(adminAccessUserList).toHaveBeenCalledWith({
			headers: {
				Authorization: 'Bearer token-123',
			},
			query: {
				order: 'desc',
				page: 3,
				pageSize: 20,
				q: 'admin',
				sort: 'email',
			},
			throwOnError: true,
		});
	});
});
