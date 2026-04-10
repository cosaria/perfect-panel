import { createFileRoute } from '@tanstack/react-router';
import { StatusCards } from '../../features/dashboard/status-cards';

export const Route = createFileRoute('/_authed/')({
	component: DashboardRoute,
});

function DashboardRoute() {
	return (
		<section className="space-y-6">
			<header className="space-y-2">
				<p className="text-xs font-medium uppercase tracking-[0.3em] text-cyan-200/80">Overview</p>
				<h1 className="text-3xl font-semibold tracking-tight text-white">Dashboard</h1>
				<p className="max-w-2xl text-sm leading-6 text-slate-300">
					这里先承载第一阶段后台总览，把认证状态、后台壳与后续工作面板统一到一个入口。
				</p>
			</header>

			<StatusCards />
		</section>
	);
}
