export const SETTINGS_DRAFT_STORAGE_KEY = 'web-v2-admin-settings-draft';

export type SettingsValues = {
	authSessionTtlMinutes: string;
	mailFromEmail: string;
	mailFromName: string;
	mailProvider: 'postmark' | 'resend' | 'smtp';
	registrationMode: 'disabled' | 'invite-only' | 'open';
	siteName: string;
	siteUrl: string;
	supportEmail: string;
};

export const defaultSettingsValues: SettingsValues = {
	authSessionTtlMinutes: '120',
	mailFromEmail: 'noreply@example.com',
	mailFromName: 'Perfect Panel',
	mailProvider: 'smtp',
	registrationMode: 'disabled',
	siteName: 'Perfect Panel',
	siteUrl: 'https://admin.example.com',
	supportEmail: 'support@example.com',
};

export function getDefaultSettingsValues(): SettingsValues {
	return {
		...defaultSettingsValues,
	};
}

export function readSettingsDraft(): SettingsValues | null {
	if (typeof window === 'undefined') {
		return null;
	}

	const raw = window.localStorage.getItem(SETTINGS_DRAFT_STORAGE_KEY);

	if (!raw) {
		return null;
	}

	try {
		const parsed = JSON.parse(raw) as Partial<SettingsValues>;

		return {
			...defaultSettingsValues,
			...parsed,
		};
	} catch {
		window.localStorage.removeItem(SETTINGS_DRAFT_STORAGE_KEY);
		return null;
	}
}

export function writeSettingsDraft(values: SettingsValues) {
	if (typeof window === 'undefined') {
		return;
	}

	window.localStorage.setItem(SETTINGS_DRAFT_STORAGE_KEY, JSON.stringify(values));
}

export function clearSettingsDraftStorage() {
	if (typeof window === 'undefined') {
		return;
	}

	window.localStorage.removeItem(SETTINGS_DRAFT_STORAGE_KEY);
}
