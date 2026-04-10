const statusCards = [
	{
		label: '认证主链',
		value: '已接通',
		description: '登录、忘记密码、重置密码已经打通到 server-v2 的真实合同。',
	},
	{
		label: '后台壳',
		value: '第一阶段',
		description: '当前先交付侧栏、顶栏和总览视图，为后续用户管理与系统设置承载内容。',
	},
	{
		label: '数据接入',
		value: '待扩展',
		description: '下一阶段再接入用户列表与系统设置骨架页，不在本任务中提前展开。',
	},
] as const;

export function StatusCards() {
	return (
		<div className="grid gap-4 md:grid-cols-3">
			{statusCards.map((card) => (
				<section
					className="rounded-3xl border border-white/10 bg-white/[0.03] p-5 shadow-2xl shadow-black/20"
					key={card.label}
				>
					<p className="text-xs font-medium uppercase tracking-[0.3em] text-slate-400">
						{card.label}
					</p>
					<div className="mt-4 text-2xl font-semibold tracking-tight text-white">{card.value}</div>
					<p className="mt-3 text-sm leading-6 text-slate-300">{card.description}</p>
				</section>
			))}
		</div>
	);
}
