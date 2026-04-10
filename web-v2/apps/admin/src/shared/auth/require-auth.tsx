import { Navigate } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import type { ReactNode } from 'react';
import { authStore } from '../store/auth-store';

export function RequireAuth({ children }: { children: ReactNode }) {
	const auth = useStore(authStore, (state) => state);

	if (!auth.accessToken) {
		return <Navigate replace to="/login" />;
	}

	return <>{children}</>;
}
