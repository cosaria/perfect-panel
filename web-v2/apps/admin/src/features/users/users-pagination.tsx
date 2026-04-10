import type { PaginationMeta } from '@web-v2/api-client';
import { Button } from '@web-v2/ui';

type UsersPaginationProps = {
	meta: PaginationMeta;
	onPageChange: (page: number) => void;
};

export function UsersPagination({ meta, onPageChange }: UsersPaginationProps) {
	return (
		<div className="flex flex-col gap-3 border-t border-white/10 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
			<p className="text-sm text-slate-400">
				共 {meta.total} 位用户 · 当前第 {meta.page} / {meta.totalPages} 页
			</p>
			<div className="flex gap-2">
				<Button
					className="rounded-xl"
					disabled={!meta.hasPrev}
					onClick={() => onPageChange(meta.page - 1)}
					variant="outline"
				>
					上一页
				</Button>
				<Button
					className="rounded-xl"
					disabled={!meta.hasNext}
					onClick={() => onPageChange(meta.page + 1)}
				>
					下一页
				</Button>
			</div>
		</div>
	);
}
