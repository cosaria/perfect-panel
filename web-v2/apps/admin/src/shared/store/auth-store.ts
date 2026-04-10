import { Store } from '@tanstack/store';

const AUTH_SESSION_STORAGE_KEY = 'web-v2-admin-auth-session';

export interface AuthState {
	accessToken: string | null;
	email: string | null;
	hydrated: boolean;
}

function readStoredAuthSession() {
	if (typeof window === 'undefined') {
		return {
			accessToken: null,
			email: null,
		};
	}

	const raw = window.sessionStorage.getItem(AUTH_SESSION_STORAGE_KEY);

	if (!raw) {
		return {
			accessToken: null,
			email: null,
		};
	}

	try {
		const parsed = JSON.parse(raw) as {
			accessToken?: string;
			email?: string;
		};

		return {
			accessToken: parsed.accessToken ?? null,
			email: parsed.email ?? null,
		};
	} catch {
		window.sessionStorage.removeItem(AUTH_SESSION_STORAGE_KEY);

		return {
			accessToken: null,
			email: null,
		};
	}
}

export const authStore = new Store<AuthState>({
	accessToken: null,
	email: null,
	hydrated: false,
});

export function hydrateAuthSession() {
	const persisted = readStoredAuthSession();

	authStore.setState(() => ({
		...persisted,
		hydrated: true,
	}));
}

export function setAuthSession(accessToken: string, email: string) {
	if (typeof window !== 'undefined') {
		window.sessionStorage.setItem(
			AUTH_SESSION_STORAGE_KEY,
			JSON.stringify({
				accessToken,
				email,
			}),
		);
	}

	authStore.setState(() => ({
		accessToken,
		email,
		hydrated: true,
	}));
}

export function clearAuthSession() {
	if (typeof window !== 'undefined') {
		window.sessionStorage.removeItem(AUTH_SESSION_STORAGE_KEY);
	}

	authStore.setState(() => ({
		accessToken: null,
		email: null,
		hydrated: true,
	}));
}
