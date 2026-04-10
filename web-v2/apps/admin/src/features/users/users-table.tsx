import { type ColumnDef, flexRender, getCoreRowModel, useReactTable } from '@tanstack/react-table';
import type { AdminAccessUserListResponse, PaginationMeta } from '@web-v2/api-client';
import type { UsersListSearch, UsersSortField } from './users-search';

type UserRow = AdminAccessUserListResponse['data'][number];

type UsersTableProps = {
	meta: PaginationMeta;
	onSortChange: (sort: UsersSortField) => void;
	rows: UserRow[];
	search: UsersListSearch;
};

const statusLabels: Record<UserRow['status'], string> = {
	active: '活跃',
	disabled: '停用',
};

export function UsersTable({ meta, onSortChange, rows, search }: UsersTableProps) {
	const renderSortMark = (field: UsersSortField) => {
		if (search.sort !== field) {
			return <span className="text-slate-500">↕</span>;
		}

		return <span className="text-cyan-200">{search.order === 'asc' ? '↑' : '↓'}</span>;
	};

	const columns: ColumnDef<UserRow>[] = [
		{
			accessorKey: 'id',
			cell: ({ getValue }) => (
				<span className="font-mono text-xs text-slate-400">{String(getValue())}</span>
			),
			header: () => (
				<button
					className="inline-flex items-center gap-2"
					onClick={() => onSortChange('id')}
					type="button"
				>
					用户 ID {renderSortMark('id')}
				</button>
			),
		},
		{
			accessorKey: 'email',
			cell: ({ getValue }) => <span className="text-sm text-white">{String(getValue())}</span>,
			header: () => (
				<button
					className="inline-flex items-center gap-2"
					onClick={() => onSortChange('email')}
					type="button"
				>
					邮箱 {renderSortMark('email')}
				</button>
			),
		},
		{
			accessorKey: 'status',
			cell: ({ getValue }) => {
				const status = getValue<UserRow['status']>();

				return (
					<span
						className={`inline-flex items-center rounded-full border px-3 py-1 text-xs font-medium ${
							status === 'active'
								? 'border-cyan-300/30 bg-cyan-300/10 text-cyan-100'
								: 'border-rose-300/30 bg-rose-300/10 text-rose-100'
						}`}
					>
						{statusLabels[status]}
					</span>
				);
			},
			header: () => (
				<button
					className="inline-flex items-center gap-2"
					onClick={() => onSortChange('status')}
					type="button"
				>
					状态 {renderSortMark('status')}
				</button>
			),
		},
	];

	const table = useReactTable({
		columns,
		data: rows,
		getCoreRowModel: getCoreRowModel(),
		manualPagination: true,
		manualSorting: true,
	});

	return (
		<div className="overflow-hidden rounded-3xl border border-white/10 bg-white/[0.03] shadow-2xl shadow-black/20">
			<div className="flex flex-col gap-3 border-b border-white/10 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
				<div>
					<h2 className="text-lg font-semibold tracking-tight text-white">真实用户列表</h2>
					<p className="text-sm leading-6 text-slate-300">
						可复制、可刷新、可按路由搜索参数恢复当前筛选状态。
					</p>
				</div>
				<div className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200">
					共 {meta.total} 位用户
				</div>
			</div>

			<div className="overflow-x-auto">
				<table className="min-w-full divide-y divide-white/10 text-left text-sm text-slate-200">
					<thead className="bg-white/[0.02] text-xs uppercase tracking-[0.25em] text-slate-400">
						{table.getHeaderGroups().map((headerGroup) => (
							<tr key={headerGroup.id}>
								{headerGroup.headers.map((header) => (
									<th className="px-5 py-4" key={header.id}>
										{header.isPlaceholder
											? null
											: flexRender(header.column.columnDef.header, header.getContext())}
									</th>
								))}
							</tr>
						))}
					</thead>
					<tbody className="divide-y divide-white/10">
						{table.getRowModel().rows.map((row) => (
							<tr className="transition hover:bg-white/[0.03]" key={row.id}>
								{row.getVisibleCells().map((cell) => (
									<td className="px-5 py-4" key={cell.id}>
										{flexRender(cell.column.columnDef.cell, cell.getContext())}
									</td>
								))}
							</tr>
						))}
					</tbody>
				</table>
			</div>

			<div className="border-t border-white/10 px-5 py-4 text-sm text-slate-400">
				第 {meta.page} / {meta.totalPages} 页 · 每页 {meta.pageSize} 条 · 当前页 {rows.length} 条
			</div>
		</div>
	);
}
