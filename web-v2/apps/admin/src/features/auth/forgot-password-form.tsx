import { useForm } from '@tanstack/react-form';
import { Link } from '@tanstack/react-router';
import { publicAuthPasswordResetRequestCreate } from '@web-v2/api-client';
import { Button } from '@web-v2/ui';
import { useState } from 'react';

type ForgotPasswordValues = {
	email: string;
};

export function ForgotPasswordForm() {
	const [message, setMessage] = useState<string | null>(null);

	const defaultValues = {
		email: '',
	} satisfies ForgotPasswordValues;

	const form = useForm({
		defaultValues: {
			...defaultValues,
		} satisfies ForgotPasswordValues,
		onSubmit: async ({ value }) => {
			setMessage(null);

			try {
				await publicAuthPasswordResetRequestCreate({
					body: value,
					throwOnError: true,
				});

				setMessage('如果该邮箱存在，我们已经发送了密码重置说明。');
			} catch {
				setMessage('提交失败，请稍后重试。');
			}
		},
	});

	return (
		<section className="min-h-screen bg-slate-950 px-6 py-10 text-slate-50">
			<div className="mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-3xl items-center justify-center">
				<div className="w-full rounded-[2rem] border border-white/10 bg-white/[0.03] p-6 shadow-2xl shadow-black/40 backdrop-blur sm:p-8">
					<div className="space-y-2">
						<h1 className="text-3xl font-semibold tracking-tight text-white">找回密码</h1>
						<p className="text-sm leading-6 text-slate-300">
							输入你的邮箱，我们会发送一封包含重置指引的邮件。
						</p>
					</div>

					<form
						className="mt-8 space-y-5"
						onSubmit={(event) => {
							event.preventDefault();
							void form.handleSubmit();
						}}
					>
						<form.Field name="email">
							{(field) => (
								<div className="space-y-2">
									<label className="block text-sm font-medium text-slate-100" htmlFor={field.name}>
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
									{isSubmitting ? '发送中' : '发送重置邮件'}
								</Button>
							)}
						</form.Subscribe>
					</form>

					<div className="mt-6 flex items-center justify-between text-sm text-slate-300">
						<Link
							className="font-medium text-cyan-200 hover:text-cyan-100"
							search={{ redirect: undefined }}
							to="/login"
						>
							返回登录
						</Link>
						<Link
							className="font-medium text-cyan-200 hover:text-cyan-100"
							search={{ token: '' }}
							to="/reset-password"
						>
							直接重置
						</Link>
					</div>
				</div>
			</div>
		</section>
	);
}
