import { Navigate } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { type ReactNode, useEffect } from 'react';
import { authStore, hydrateAuthSession } from '../store/auth-store';

export function RequireAuth({ children }: { children: ReactNode }) {
	const auth = useStore(authStore, (state) => state);

	useEffect(() => {
		if (!auth.hydrated) {
			hydrateAuthSession();
		}
	}, [auth.hydrated]);

	if (!auth.hydrated) {
		return null;
	}

	if (!auth.accessToken) {
		return <Navigate replace to="/login" />;
	}

	return <>{children}</>;
}
