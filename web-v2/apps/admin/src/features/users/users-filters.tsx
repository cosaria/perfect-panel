import { Button } from '@web-v2/ui';
import type { ChangeEvent, FormEvent } from 'react';
import { DEFAULT_USERS_PAGE_SIZE, type UsersListSearch } from './users-search';

type UsersFiltersProps = {
	isRefreshing: boolean;
	onPageSizeChange: (pageSize: number) => void;
	onQuerySubmit: () => void;
	onRefresh: () => void;
	onSearchTextChange: (value: string) => void;
	pageSizeOptions?: readonly number[];
	search: UsersListSearch;
	searchText: string;
};

const defaultPageSizes = [10, 20, 50] as const;

export function UsersFilters({
	isRefreshing,
	onPageSizeChange,
	onQuerySubmit,
	onRefresh,
	onSearchTextChange,
	pageSizeOptions = defaultPageSizes,
	searchText,
	search,
}: UsersFiltersProps) {
	const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		onQuerySubmit();
	};

	const handlePageSizeChange = (event: ChangeEvent<HTMLSelectElement>) => {
		onPageSizeChange(Number(event.target.value));
	};

	return (
		<div className="rounded-3xl border border-white/10 bg-white/[0.03] p-4 shadow-2xl shadow-black/20">
			<form
				className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_12rem_auto_auto]"
				onSubmit={handleSubmit}
			>
				<div className="space-y-2">
					<label
						className="block text-xs font-medium uppercase tracking-[0.28em] text-slate-400"
						htmlFor="users-search"
					>
						搜索用户
					</label>
					<input
						autoComplete="off"
						className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition placeholder:text-slate-500 focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
						id="users-search"
						name="users-search"
						onChange={(event) => onSearchTextChange(event.target.value)}
						placeholder="搜索邮箱或状态"
						value={searchText}
					/>
				</div>

				<div className="space-y-2">
					<label
						className="block text-xs font-medium uppercase tracking-[0.28em] text-slate-400"
						htmlFor="users-page-size"
					>
						每页条数
					</label>
					<select
						className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
						id="users-page-size"
						name="users-page-size"
						onChange={handlePageSizeChange}
						value={search.pageSize}
					>
						{pageSizeOptions.map((value) => (
							<option key={value} value={value}>
								{value}
							</option>
						))}
					</select>
				</div>

				<div className="flex items-end">
					<Button className="h-11 w-full rounded-xl text-sm lg:w-auto" type="submit">
						查询
					</Button>
				</div>

				<div className="flex items-end">
					<Button
						className="h-11 w-full rounded-xl text-sm lg:w-auto"
						disabled={isRefreshing}
						onClick={onRefresh}
						type="button"
						variant="outline"
					>
						{isRefreshing ? '刷新中' : '刷新'}
					</Button>
				</div>
			</form>

			<p className="mt-3 text-xs leading-5 text-slate-400">
				当前列表状态来自路由搜索参数，默认每页 {DEFAULT_USERS_PAGE_SIZE} 条。
			</p>
		</div>
	);
}
