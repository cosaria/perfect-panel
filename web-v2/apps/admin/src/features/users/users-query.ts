import { keepPreviousData, useQuery } from '@tanstack/react-query';
import { type AdminAccessUserListResponse, adminAccessUserList } from '@web-v2/api-client';
import { type UsersListSearch, buildUsersListQuery } from './users-search';

export function createUsersListQueryKey(search: UsersListSearch) {
	return ['admin', 'users', buildUsersListQuery(search)] as const;
}

export async function fetchUsersList(search: UsersListSearch, accessToken: string) {
	const response = await adminAccessUserList({
		headers: {
			Authorization: `Bearer ${accessToken}`,
		},
		query: buildUsersListQuery(search),
		throwOnError: true,
	});

	return response.data;
}

export function useUsersListQuery(search: UsersListSearch, accessToken: string | null) {
	return useQuery<AdminAccessUserListResponse>({
		enabled: Boolean(accessToken),
		placeholderData: keepPreviousData,
		queryFn: () => {
			if (!accessToken) {
				throw new Error('缺少管理员会话令牌。');
			}

			return fetchUsersList(search, accessToken);
		},
		queryKey: [...createUsersListQueryKey(search), accessToken] as const,
	});
}
