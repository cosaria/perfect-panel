import { Link } from '@tanstack/react-router';
import { Sidebar } from '@web-v2/ui';

type NavigationItem = {
	label: string;
	description: string;
	to?: '/';
};

const navigationItems: NavigationItem[] = [
	{
		label: 'Dashboard',
		description: '查看第一阶段的总览状态。',
		to: '/',
	},
	{
		label: '用户管理',
		description: '下一阶段接入用户列表与筛选。',
	},
	{
		label: '系统设置',
		description: '下一阶段接入系统配置表单。',
	},
];

export function AdminSidebar() {
	return (
		<Sidebar className="w-full border-white/10 bg-slate-950/95 text-slate-50 lg:w-72">
			<div className="border-b border-white/10 px-5 py-6">
				<p className="text-xs font-medium uppercase tracking-[0.35em] text-cyan-200/80">
					Perfect Panel
				</p>
				<div className="mt-3 space-y-2">
					<h1 className="text-xl font-semibold tracking-tight text-white">Admin Console</h1>
					<p className="text-sm leading-6 text-slate-300">
						基于第一阶段导航骨架，先把后台壳与总览体验搭起来。
					</p>
				</div>
			</div>

			<nav className="flex-1 space-y-2 px-3 py-4">
				{navigationItems.map((item) =>
					item.to ? (
						<Link
							activeOptions={{ exact: item.to === '/' }}
							activeProps={{
								className: 'border-cyan-300/40 bg-cyan-300/10 text-cyan-100',
							}}
							className="block rounded-2xl border border-transparent px-4 py-3 transition hover:border-white/10 hover:bg-white/5"
							key={item.label}
							to={item.to}
						>
							<div className="text-sm font-medium">{item.label}</div>
							<p className="mt-1 text-xs leading-5 text-slate-400">{item.description}</p>
						</Link>
					) : (
						<div
							aria-disabled="true"
							className="rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 text-slate-400"
							key={item.label}
						>
							<div className="text-sm font-medium text-slate-200">{item.label}</div>
							<p className="mt-1 text-xs leading-5 text-slate-400">{item.description}</p>
						</div>
					),
				)}
			</nav>
		</Sidebar>
	);
}
