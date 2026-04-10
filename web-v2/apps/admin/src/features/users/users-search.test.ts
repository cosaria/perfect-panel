import { describe, expect, it } from 'vitest';
import {
	DEFAULT_USERS_ORDER,
	DEFAULT_USERS_PAGE,
	DEFAULT_USERS_PAGE_SIZE,
	DEFAULT_USERS_SORT,
	buildUsersListQuery,
	parseUsersSearch,
	toggleUsersSort,
	updateUsersSearch,
} from './users-search';

describe('users search helpers', () => {
	it('normalizes invalid search values back to safe defaults', () => {
		expect(
			parseUsersSearch({
				order: 'sideways',
				page: '-4',
				pageSize: '999',
				q: 42,
				sort: 'createdAt',
			}),
		).toEqual({
			order: DEFAULT_USERS_ORDER,
			page: DEFAULT_USERS_PAGE,
			pageSize: DEFAULT_USERS_PAGE_SIZE,
			q: '',
			sort: DEFAULT_USERS_SORT,
		});
	});

	it('builds list query values with empty q stripped out', () => {
		expect(
			buildUsersListQuery({
				order: 'desc',
				page: 2,
				pageSize: 20,
				q: '   ',
				sort: 'status',
			}),
		).toEqual({
			order: 'desc',
			page: 2,
			pageSize: 20,
			q: undefined,
			sort: 'status',
		});
	});

	it('toggles sort order and resets page', () => {
		expect(
			toggleUsersSort(
				{
					order: 'asc',
					page: 4,
					pageSize: 10,
					q: 'admin',
					sort: 'email',
				},
				'email',
			),
		).toEqual({
			order: 'desc',
			page: 1,
			pageSize: 10,
			q: 'admin',
			sort: 'email',
		});
	});

	it('updates search while preserving normalized state', () => {
		expect(
			updateUsersSearch(
				{
					order: 'asc',
					page: 3,
					pageSize: 10,
					q: '',
					sort: 'email',
				},
				{
					page: 0,
					pageSize: 50,
					q: 'hello',
				},
			),
		).toEqual({
			order: 'asc',
			page: 1,
			pageSize: 50,
			q: 'hello',
			sort: 'email',
		});
	});
});
