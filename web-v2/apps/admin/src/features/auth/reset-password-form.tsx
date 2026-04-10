import { useForm } from '@tanstack/react-form';
import { Link, useNavigate } from '@tanstack/react-router';
import { publicAuthPasswordResetCreate } from '@web-v2/api-client';
import { Button } from '@web-v2/ui';
import { useState } from 'react';

type ResetPasswordValues = {
	token: string;
	password: string;
};

export function ResetPasswordForm({ initialToken = '' }: { initialToken?: string }) {
	const navigate = useNavigate();
	const [message, setMessage] = useState<string | null>(null);

	const defaultValues = {
		token: initialToken,
		password: '',
	} satisfies ResetPasswordValues;

	const form = useForm({
		defaultValues: {
			...defaultValues,
		} satisfies ResetPasswordValues,
		onSubmit: async ({ value }) => {
			setMessage(null);

			try {
				await publicAuthPasswordResetCreate({
					body: value,
					throwOnError: true,
				});

				setMessage('密码已经更新，现在可以重新登录。');
				await navigate({ replace: true, to: '/login' });
			} catch {
				setMessage('重置失败，请检查令牌是否正确。');
			}
		},
	});

	return (
		<section className="min-h-screen bg-slate-950 px-6 py-10 text-slate-50">
			<div className="mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-3xl items-center justify-center">
				<div className="w-full rounded-[2rem] border border-white/10 bg-white/[0.03] p-6 shadow-2xl shadow-black/40 backdrop-blur sm:p-8">
					<div className="space-y-2">
						<h1 className="text-3xl font-semibold tracking-tight text-white">重置密码</h1>
						<p className="text-sm leading-6 text-slate-300">输入令牌和新密码，完成账户密码更新。</p>
					</div>

					<form
						className="mt-8 space-y-5"
						onSubmit={(event) => {
							event.preventDefault();
							void form.handleSubmit();
						}}
					>
						<form.Field name="token">
							{(field) => (
								<div className="space-y-2">
									<label className="block text-sm font-medium text-slate-100" htmlFor={field.name}>
										令牌
									</label>
									<input
										autoComplete="one-time-code"
										className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition placeholder:text-slate-500 focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
										id={field.name}
										name={field.name}
										onBlur={field.handleBlur}
										onChange={(event) => field.handleChange(event.target.value)}
										placeholder="请输入重置令牌"
										type="text"
										value={field.state.value}
									/>
								</div>
							)}
						</form.Field>

						<form.Field name="password">
							{(field) => (
								<div className="space-y-2">
									<label className="block text-sm font-medium text-slate-100" htmlFor={field.name}>
										新密码
									</label>
									<input
										autoComplete="new-password"
										className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition placeholder:text-slate-500 focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
										id={field.name}
										name={field.name}
										onBlur={field.handleBlur}
										onChange={(event) => field.handleChange(event.target.value)}
										placeholder="请输入新密码"
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
									{isSubmitting ? '提交中' : '确认重置'}
								</Button>
							)}
						</form.Subscribe>
					</form>

					<div className="mt-6 flex items-center justify-between text-sm text-slate-300">
						<Link className="font-medium text-cyan-200 hover:text-cyan-100" to="/login">
							返回登录
						</Link>
						<Link className="font-medium text-cyan-200 hover:text-cyan-100" to="/forgot-password">
							重新发送邮件
						</Link>
					</div>
				</div>
			</div>
		</section>
	);
}
