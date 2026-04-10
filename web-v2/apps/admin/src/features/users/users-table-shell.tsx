export function UsersTableShell() {
	return (
		<section className="space-y-6">
			<header className="flex items-center justify-between">
				<div>
					<h1 className="text-2xl font-semibold tracking-tight text-white">用户管理</h1>
					<p className="text-sm leading-6 text-slate-300">当前阶段先提供页面结构和空状态。</p>
				</div>
			</header>

			<div className="rounded-3xl border border-white/10 bg-white/[0.03] p-6 text-sm leading-6 text-slate-300">
				暂无用户数据。后续任务会在这里接入真实列表和筛选。
			</div>
		</section>
	);
}
