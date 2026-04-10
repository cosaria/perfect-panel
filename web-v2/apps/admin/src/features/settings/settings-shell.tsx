export function SettingsShell() {
	return (
		<section className="space-y-6">
			<header className="flex items-center justify-between">
				<div>
					<h1 className="text-2xl font-semibold tracking-tight text-white">系统设置</h1>
					<p className="text-sm leading-6 text-slate-300">
						当前阶段先提供页面骨架，后续任务再接入真实配置项。
					</p>
				</div>
			</header>

			<div className="rounded-3xl border border-white/10 bg-white/[0.03] p-6 text-sm leading-6 text-slate-300">
				系统配置表单尚未接入。下一阶段会在这里落邮件、认证与站点设置项。
			</div>
		</section>
	);
}
