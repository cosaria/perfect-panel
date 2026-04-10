import { useNavigate } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { useEffect, useState } from 'react';
import { authStore } from '../../shared/store/auth-store';
import { UsersFilters } from './users-filters';
import { UsersPagination } from './users-pagination';
import { useUsersListQuery } from './users-query';
import {
	DEFAULT_USERS_ORDER,
	DEFAULT_USERS_PAGE,
	DEFAULT_USERS_PAGE_SIZE,
	DEFAULT_USERS_SORT,
	type UsersListSearch,
	type UsersSortField,
	toggleUsersSort,
	updateUsersSearch,
} from './users-search';
import { UsersTable } from './users-table';

type UsersTableShellProps = {
	search: UsersListSearch;
};

export function UsersTableShell({ search }: UsersTableShellProps) {
	const navigate = useNavigate();
	const accessToken = useStore(authStore, (state) => state.accessToken);
	const query = useUsersListQuery(search, accessToken);
	const [searchText, setSearchText] = useState(search.q);

	useEffect(() => {
		setSearchText(search.q);
	}, [search.q]);

	const updateSearch = (patch: Partial<UsersListSearch>) => {
		const nextSearch = updateUsersSearch(search, patch);

		void navigate({
			search: nextSearch,
			to: '/users',
		});
	};

	const handleSortChange = (sort: UsersSortField) => {
		const nextSearch = toggleUsersSort(search, sort);

		void navigate({
			search: nextSearch,
			to: '/users',
		});
	};

	const handleQuerySubmit = () => {
		updateSearch({
			page: DEFAULT_USERS_PAGE,
			pageSize: search.pageSize || DEFAULT_USERS_PAGE_SIZE,
			q: searchText.trim(),
			sort: search.sort || DEFAULT_USERS_SORT,
			order: search.order || DEFAULT_USERS_ORDER,
		});
	};

	const handlePageSizeChange = (pageSize: number) => {
		updateSearch({
			page: DEFAULT_USERS_PAGE,
			pageSize,
		});
	};

	const handlePageChange = (page: number) => {
		updateSearch({
			page,
		});
	};

	const handleRefresh = () => {
		void query.refetch();
	};

	const data = query.data?.data ?? [];
	const meta = query.data?.meta;
	const isEmpty = query.isSuccess && data.length === 0;

	return (
		<section className="space-y-6">
			<header className="space-y-4">
				<div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
					<div className="space-y-2">
						<p className="text-xs font-medium uppercase tracking-[0.35em] text-cyan-200/80">
							Admin / Users
						</p>
						<h1 className="text-2xl font-semibold tracking-tight text-white">用户管理</h1>
						<p className="max-w-2xl text-sm leading-6 text-slate-300">
							真实用户列表直接由路由搜索参数驱动，支持搜索、分页和排序，刷新后也能恢复当前筛选状态。
						</p>
					</div>

					<div className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200">
						共 {query.data?.meta.total ?? 0} 位用户
					</div>
				</div>

				<UsersFilters
					isRefreshing={query.isFetching}
					onPageSizeChange={handlePageSizeChange}
					onQuerySubmit={handleQuerySubmit}
					onRefresh={handleRefresh}
					onSearchTextChange={setSearchText}
					pageSizeOptions={[10, 20, 50]}
					search={search}
					searchText={searchText}
				/>
			</header>

			{query.isLoading ? (
				<div className="rounded-3xl border border-white/10 bg-white/[0.03] p-5 shadow-2xl shadow-black/20">
					<div className="space-y-3">
						<div className="h-5 w-40 rounded-full bg-white/10" />
						<div className="grid gap-3 pt-3">
							{['s1', 's2', 's3', 's4', 's5'].map((key) => (
								<div
									className="h-14 rounded-2xl border border-white/5 bg-white/[0.03] animate-pulse"
									key={key}
								/>
							))}
						</div>
					</div>
				</div>
			) : query.isError ? (
				<div className="rounded-3xl border border-rose-300/20 bg-rose-300/10 p-6 text-sm leading-6 text-rose-100">
					<div className="space-y-2">
						<h2 className="text-base font-semibold text-rose-50">用户列表加载失败</h2>
						<p>{query.error instanceof Error ? query.error.message : '请稍后重试。'}</p>
					</div>
					<div className="mt-4">
						<button
							className="rounded-xl border border-rose-200/30 px-4 py-2 text-sm text-rose-50 transition hover:bg-rose-200/10"
							onClick={handleRefresh}
							type="button"
						>
							重新加载
						</button>
					</div>
				</div>
			) : isEmpty ? (
				<div className="rounded-3xl border border-white/10 bg-white/[0.03] p-8 text-sm leading-6 text-slate-300 shadow-2xl shadow-black/20">
					<div className="max-w-xl space-y-3">
						<h2 className="text-lg font-semibold tracking-tight text-white">没有找到匹配的用户</h2>
						<p>
							当前筛选条件没有返回结果。可以尝试清空搜索词、切换排序方式，或者调整每页条数后再次查询。
						</p>
					</div>
				</div>
			) : meta ? (
				<div className="space-y-4">
					{query.isFetching ? (
						<div className="rounded-full border border-cyan-300/20 bg-cyan-300/10 px-4 py-2 text-xs font-medium uppercase tracking-[0.3em] text-cyan-100">
							正在刷新结果
						</div>
					) : null}
					<UsersTable meta={meta} onSortChange={handleSortChange} rows={data} search={search} />
					<UsersPagination meta={meta} onPageChange={handlePageChange} />
				</div>
			) : null}
		</section>
	);
}
