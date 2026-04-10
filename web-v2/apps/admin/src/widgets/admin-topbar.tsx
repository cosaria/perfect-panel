import { useStore } from '@tanstack/react-store';
import { authStore } from '../shared/store/auth-store';

export function AdminTopbar() {
	const auth = useStore(authStore, (state) => state);

	return (
		<header className="border-b border-white/10 bg-slate-950/70 px-6 py-4 backdrop-blur lg:px-8">
			<div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
				<div className="space-y-1">
					<p className="text-xs font-medium uppercase tracking-[0.3em] text-slate-400">
						Admin Shell
					</p>
					<h2 className="text-lg font-semibold tracking-tight text-white">第一阶段后台骨架</h2>
				</div>
				<div className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200">
					{auth.email ?? '未识别管理员'}
				</div>
			</div>
		</header>
	);
}
