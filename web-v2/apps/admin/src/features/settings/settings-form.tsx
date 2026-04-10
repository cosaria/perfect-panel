import { useForm } from '@tanstack/react-form';
import { Button } from '@web-v2/ui';
import { useState } from 'react';
import {
	type SettingsValues,
	clearSettingsDraftStorage,
	getDefaultSettingsValues,
	readSettingsDraft,
	writeSettingsDraft,
} from './settings-draft';
import {
	type SettingsFieldName,
	type SettingsValidationErrors,
	getFieldError,
	validateSettingsValues,
} from './settings-schema';

const fieldGroups: Array<{
	description: string;
	fields: Array<{
		label: string;
		name: SettingsFieldName;
		placeholder?: string;
		type?: 'email' | 'number' | 'text' | 'url';
	}>;
	title: string;
}> = [
	{
		description: '这些字段决定后台和通知里展示出来的站点基本信息。',
		fields: [
			{
				label: '站点名称',
				name: 'siteName',
				placeholder: 'Perfect Panel',
			},
			{
				label: '站点地址',
				name: 'siteUrl',
				placeholder: 'https://admin.example.com',
				type: 'url',
			},
			{
				label: '支持邮箱',
				name: 'supportEmail',
				placeholder: 'support@example.com',
				type: 'email',
			},
		],
		title: '站点设置',
	},
	{
		description: '这里只先做前端草稿，字段名称会贴近后续 system_settings 的真实落点。',
		fields: [
			{
				label: '会话时长（分钟）',
				name: 'authSessionTtlMinutes',
				placeholder: '120',
				type: 'number',
			},
		],
		title: '认证设置',
	},
	{
		description: '先把邮件发送的品牌和通道配置做成可编辑草稿。',
		fields: [
			{
				label: '发件人名称',
				name: 'mailFromName',
				placeholder: 'Perfect Panel',
			},
			{
				label: '发件邮箱',
				name: 'mailFromEmail',
				placeholder: 'noreply@example.com',
				type: 'email',
			},
		],
		title: '邮件设置',
	},
];

export function SettingsForm() {
	const initialDraft = readSettingsDraft();
	const [statusMessage, setStatusMessage] = useState<string | null>(
		initialDraft ? '已从本地恢复上次保存的设置草稿。' : null,
	);
	const [validationErrors, setValidationErrors] = useState<SettingsValidationErrors>({});

	const form = useForm({
		defaultValues: initialDraft ?? getDefaultSettingsValues(),
	});

	const clearFieldError = (fieldName: SettingsFieldName) => {
		setValidationErrors((current) => {
			if (!current[fieldName]) {
				return current;
			}

			const next = { ...current };
			delete next[fieldName];
			return next;
		});
	};

	const handleClearDraft = () => {
		clearSettingsDraftStorage();
		setValidationErrors({});
		form.reset(getDefaultSettingsValues());
		setStatusMessage('本地草稿已清空，表单已恢复默认值。');
	};

	return (
		<section className="space-y-6">
			<header className="space-y-3">
				<div className="space-y-2">
					<p className="text-xs font-medium uppercase tracking-[0.35em] text-cyan-200/80">
						Admin / Settings
					</p>
					<h1 className="text-2xl font-semibold tracking-tight text-white">系统设置</h1>
					<p className="max-w-2xl text-sm leading-6 text-slate-300">
						这轮先把设置页做成真实前端表单和本地草稿体验，等 server-v2 暴露 settings API
						后再无缝接上。
					</p>
				</div>
				<div className="rounded-2xl border border-amber-300/20 bg-amber-300/10 px-4 py-3 text-sm text-amber-50">
					当前为前端草稿模式，保存后只会写入本地浏览器，不会提交到服务端。
				</div>
			</header>

			<form
				className="space-y-6"
				noValidate
				onSubmit={(event) => {
					event.preventDefault();
					const values = form.state.values;
					const errors = validateSettingsValues(values);
					setValidationErrors(errors);

					if (Object.keys(errors).length > 0) {
						setStatusMessage('请先修正表单错误，再保存草稿。');
						return;
					}

					writeSettingsDraft(values);
					setStatusMessage('本地草稿已保存，当前尚未同步到服务端。');
				}}
			>
				{fieldGroups.map((group) => (
					<div
						className="rounded-3xl border border-white/10 bg-white/[0.03] p-6 shadow-2xl shadow-black/20"
						key={group.title}
					>
						<div className="mb-5 space-y-2">
							<h2 className="text-lg font-semibold tracking-tight text-white">{group.title}</h2>
							<p className="text-sm leading-6 text-slate-300">{group.description}</p>
						</div>

						<div className="grid gap-5 md:grid-cols-2">
							{group.fields.map((fieldConfig) => (
								<div className="space-y-2" key={fieldConfig.name}>
									<form.Field name={fieldConfig.name}>
										{(field) => (
											<>
												<label
													className="block text-sm font-medium text-slate-100"
													htmlFor={field.name}
												>
													{fieldConfig.label}
												</label>
												<input
													className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition placeholder:text-slate-500 focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
													id={field.name}
													name={field.name}
													onBlur={field.handleBlur}
													onChange={(event) => {
														clearFieldError(fieldConfig.name);
														field.handleChange(event.target.value);
													}}
													placeholder={fieldConfig.placeholder}
													type={fieldConfig.type ?? 'text'}
													value={field.state.value}
												/>
											</>
										)}
									</form.Field>
									{getFieldError(validationErrors, fieldConfig.name) ? (
										<p className="text-sm text-rose-200">
											{getFieldError(validationErrors, fieldConfig.name)}
										</p>
									) : null}
								</div>
							))}

							{group.title === '认证设置' ? (
								<form.Field name="registrationMode">
									{(field) => (
										<div className="space-y-2">
											<label
												className="block text-sm font-medium text-slate-100"
												htmlFor={field.name}
											>
												注册模式
											</label>
											<select
												className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
												id={field.name}
												name={field.name}
												onBlur={field.handleBlur}
												onChange={(event) =>
													field.handleChange(
														event.target.value as SettingsValues['registrationMode'],
													)
												}
												value={field.state.value}
											>
												<option value="disabled">关闭注册</option>
												<option value="invite-only">仅邀请</option>
												<option value="open">公开注册</option>
											</select>
										</div>
									)}
								</form.Field>
							) : null}

							{group.title === '邮件设置' ? (
								<form.Field name="mailProvider">
									{(field) => (
										<div className="space-y-2">
											<label
												className="block text-sm font-medium text-slate-100"
												htmlFor={field.name}
											>
												邮件通道
											</label>
											<select
												className="h-11 w-full rounded-xl border border-white/10 bg-slate-900/80 px-4 text-sm text-white outline-none ring-0 transition focus:border-cyan-300/60 focus:ring-2 focus:ring-cyan-300/20"
												id={field.name}
												name={field.name}
												onBlur={field.handleBlur}
												onChange={(event) =>
													field.handleChange(event.target.value as SettingsValues['mailProvider'])
												}
												value={field.state.value}
											>
												<option value="smtp">SMTP</option>
												<option value="resend">Resend</option>
												<option value="postmark">Postmark</option>
											</select>
										</div>
									)}
								</form.Field>
							) : null}
						</div>
					</div>
				))}

				{statusMessage ? (
					<div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-200">
						{statusMessage}
					</div>
				) : null}

				<form.Subscribe
					selector={(state) => ({
						isDirty: state.isDirty,
						isSubmitting: state.isSubmitting,
					})}
				>
					{({ isDirty, isSubmitting }) => (
						<div className="flex flex-col gap-3 sm:flex-row">
							<Button className="rounded-xl" disabled={isSubmitting} type="submit">
								{isSubmitting ? '保存中' : '保存草稿'}
							</Button>
							<Button
								className="rounded-xl"
								onClick={handleClearDraft}
								type="button"
								variant="outline"
							>
								清空草稿
							</Button>
							<div className="flex items-center rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-300">
								{isDirty ? '当前表单有未保存改动' : '当前表单与最近一次加载状态一致'}
							</div>
						</div>
					)}
				</form.Subscribe>
			</form>
		</section>
	);
}
