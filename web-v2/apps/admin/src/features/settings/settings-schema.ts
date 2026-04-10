import type { SettingsValues } from './settings-draft';

export type SettingsFieldName = keyof SettingsValues;
export type SettingsValidationErrors = Partial<Record<SettingsFieldName, string>>;

function isEmail(value: string) {
	return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
}

function isPositiveInteger(value: string) {
	return Number.isInteger(Number(value)) && Number(value) > 0;
}

function isUrl(value: string) {
	try {
		const url = new URL(value);
		return url.protocol === 'http:' || url.protocol === 'https:';
	} catch {
		return false;
	}
}

export function validateSettingsValues(values: SettingsValues): SettingsValidationErrors {
	const errors: SettingsValidationErrors = {};

	if (!values.siteName.trim()) {
		errors.siteName = '请输入站点名称。';
	}

	if (!isUrl(values.siteUrl.trim())) {
		errors.siteUrl = '请输入有效的站点地址。';
	}

	if (!isEmail(values.supportEmail.trim())) {
		errors.supportEmail = '请输入有效的支持邮箱。';
	}

	if (!isPositiveInteger(values.authSessionTtlMinutes.trim())) {
		errors.authSessionTtlMinutes = '会话时长必须是正整数分钟。';
	}

	if (!values.mailFromName.trim()) {
		errors.mailFromName = '请输入发件人名称。';
	}

	if (!isEmail(values.mailFromEmail.trim())) {
		errors.mailFromEmail = '请输入有效的发件邮箱。';
	}

	return errors;
}

export function getFieldError(errors: SettingsValidationErrors, field: SettingsFieldName) {
	return errors[field];
}
