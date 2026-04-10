import type { Order as SortOrder } from '@web-v2/api-client';

export const DEFAULT_USERS_PAGE = 1;
export const DEFAULT_USERS_PAGE_SIZE = 10;
export const DEFAULT_USERS_SORT = 'email';
export const DEFAULT_USERS_ORDER: SortOrder = 'asc';

export type UsersSortField = 'id' | 'email' | 'status';

export type UsersListSearch = {
	q: string;
	page: number;
	pageSize: number;
	sort: UsersSortField;
	order: SortOrder;
};

const USERS_PAGE_SIZES = [10, 20, 50] as const;

function parsePositiveInteger(value: unknown, fallback: number) {
	if (typeof value === 'number' && Number.isInteger(value) && value > 0) {
		return value;
	}

	if (typeof value === 'string') {
		const parsed = Number.parseInt(value, 10);
		if (Number.isInteger(parsed) && parsed > 0) {
			return parsed;
		}
	}

	return fallback;
}

function parseSortField(value: unknown): UsersSortField {
	if (value === 'id' || value === 'email' || value === 'status') {
		return value;
	}

	return DEFAULT_USERS_SORT;
}

function parseOrder(value: unknown): SortOrder {
	return value === 'desc' ? 'desc' : DEFAULT_USERS_ORDER;
}

export function parseUsersSearch(search: Record<string, unknown>): UsersListSearch {
	const pageSize = parsePositiveInteger(search.pageSize, DEFAULT_USERS_PAGE_SIZE);

	return {
		q: typeof search.q === 'string' ? search.q : '',
		page: parsePositiveInteger(search.page, DEFAULT_USERS_PAGE),
		pageSize: USERS_PAGE_SIZES.includes(pageSize as (typeof USERS_PAGE_SIZES)[number])
			? pageSize
			: DEFAULT_USERS_PAGE_SIZE,
		sort: parseSortField(search.sort),
		order: parseOrder(search.order),
	};
}

export function buildUsersListQuery(search: UsersListSearch) {
	return {
		order: search.order,
		page: search.page,
		pageSize: search.pageSize,
		q: search.q.trim() || undefined,
		sort: search.sort,
	};
}

export function normalizeUsersSearch(search: UsersListSearch): UsersListSearch {
	return parseUsersSearch(search);
}

export function updateUsersSearch(
	current: UsersListSearch,
	patch: Partial<UsersListSearch>,
): UsersListSearch {
	return normalizeUsersSearch({
		...current,
		...patch,
	});
}

export function toggleUsersSort(current: UsersListSearch, sort: UsersSortField): UsersListSearch {
	const nextOrder: SortOrder = current.sort === sort && current.order === 'asc' ? 'desc' : 'asc';

	return updateUsersSearch(current, {
		order: nextOrder,
		page: DEFAULT_USERS_PAGE,
		sort,
	});
}
