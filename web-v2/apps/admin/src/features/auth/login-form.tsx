import { useForm } from '@tanstack/react-form';
import { Link, useNavigate } from '@tanstack/react-router';
import { publicAuthSessionCreate } from '@web-v2/api-client';
import { Button } from '@web-v2/ui';
import { useState } from 'react';
import { setAuthSession } from '../../shared/store/auth-store';

type LoginValues = {
	email: string;
	password: string;
};

export function LoginForm() {
	const navigate = useNavigate();
	const [message, setMessage] = useState<string | null>(null);

	const defaultValues = {
		email: '',
		password: '',
	} satisfies LoginValues;

	const form = useForm({
		defaultValues: {
			...defaultValues,
		} satisfies LoginValues,
		onSubmit: async ({ value }) => {
			setMessage(null);

			try {
				const response = await publicAuthSessionCreate({
					body: value,
					throwOnError: true,
				});

				setAuthSession(response.data.data.accessToken, value.email);
				setMessage('登录成功，正在进入管理后台。');
				await navigate({ replace: true, to: '/' });
			} catch {
				setMessage('登录失败，请检查邮箱和密码后重试。');
			}
		},
	});

	return (
		<section className="min-h-screen bg-slate-950 px-6 py-10 text-slate-50">
			<div className="mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-6xl items-center justify-center">
				<div className="grid w-full overflow-hidden rounded-[2rem] border border-white/10 bg-white/[0.03] shadow-2xl shadow-black/40 backdrop-blur md:grid-cols-[1.1fr_0.9fr]">
					<div className="hidden flex-col justify-between border-r border-white/10 p-10 md:flex">
						<div className="space-y-4">
							<p className="text-sm font-medium uppercase tracking-[0.35em] text-cyan-200/80">
								Perfect Panel
							</p>
							<h1 className="max-w-md text-4xl font-semibold tracking-tight text-white">
								回到管理后台，继续处理你的业务。
							</h1>
							<p className="max-w-md text-sm leading-6 text-slate-300">
								统一登录、找回密码与后续受保护页面会在这里衔接，先把认证主链路打通。
							</p>
						</div>
						<div className="rounded-2xl border border-white/10 bg-white/5 p-4 text-sm text-slate-200">
							<div className="mb-2 text-xs uppercase tracking-[0.3em] text-slate-400">状态</div>
							<p>当前只开放认证入口，其他管理壳在后续任务接入。</p>
						</div>
					</div>
					<div className="p-6 sm:p-8 lg:p-10">
						<div className="mx-auto flex w-full max-w-md flex-col gap-8">
							<div className="space-y-3">
								<div className="inline-flex items-center rounded-full border border-cyan-300/20 bg-cyan-300/10 px-3 py-1 text-xs font-medium text-cyan-100">
									管理员登录
								</div>
								<div className="space-y-2">
									<h2 className="text-3xl font-semibold tracking-tight text-white">登录你的账户</h2>
									<p className="text-sm leading-6 text-slate-300">
										请输入邮箱和密码，继续进入后台。
									</p>
								</div>
							</div>

							<form
								className="space-y-5"
								onSubmit={(event) => {
									event.preventDefault();
									void form.handleSubmit();
								}}
							>
								<form.Field name="email">
									{(field) => (
										<div className="space-y-2">
											<label
												className="block text-sm font-medium text-slate-100"
												htmlFor={field.name}
											>
												邮箱
											</label>
											<input
												autoComplete="email"
												className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition placeholder:text-slate-500 focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
												id={field.name}
												name={field.name}
												onBlur={field.handleBlur}
												onChange={(event) => field.handleChange(event.target.value)}
												placeholder="admin@example.com"
												type="email"
												value={field.state.value}
											/>
										</div>
									)}
								</form.Field>

								<form.Field name="password">
									{(field) => (
										<div className="space-y-2">
											<label
												className="block text-sm font-medium text-slate-100"
												htmlFor={field.name}
											>
												密码
											</label>
											<input
												autoComplete="current-password"
												className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition placeholder:text-slate-500 focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
												id={field.name}
												name={field.name}
												onBlur={field.handleBlur}
												onChange={(event) => field.handleChange(event.target.value)}
												placeholder="请输入密码"
												type="password"
												value={field.state.value}
											/>
										</div>
									)}
								</form.Field>

								{message ? (
									<p className="rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-200">
										{message}
									</p>
								) : null}

								<form.Subscribe
									selector={(state) => ({
										canSubmit: state.canSubmit,
										isSubmitting: state.isSubmitting,
									})}
								>
									{({ canSubmit, isSubmitting }) => (
										<Button
											className="h-11 w-full rounded-xl text-sm"
											disabled={!canSubmit || isSubmitting}
											type="submit"
										>
											{isSubmitting ? '登录中' : '登录'}
										</Button>
									)}
								</form.Subscribe>
							</form>

							<div className="flex items-center justify-between text-sm text-slate-300">
								<span>忘记密码了？</span>
								<Link
									className="font-medium text-cyan-200 hover:text-cyan-100"
									to="/forgot-password"
								>
									前往找回
								</Link>
							</div>
						</div>
					</div>
				</div>
			</div>
		</section>
	);
}
