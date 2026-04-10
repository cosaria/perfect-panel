import { Store } from '@tanstack/store';

export interface AuthState {
	accessToken: string | null;
	email: string | null;
}

export const authStore = new Store<AuthState>({
	accessToken: null,
	email: null,
});

export function setAuthSession(accessToken: string, email: string) {
	authStore.setState(() => ({
		accessToken,
		email,
	}));
}

export function clearAuthSession() {
	authStore.setState(() => ({
		accessToken: null,
		email: null,
	}));
}
