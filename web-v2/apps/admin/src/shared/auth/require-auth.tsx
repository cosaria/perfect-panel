import { Navigate, useRouterState } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { type ReactNode, useEffect, useRef } from 'react';
import { authStore, hydrateAuthSession } from '../store/auth-store';
import { buildRedirectTarget } from './redirect';

export function AuthHydrationFallback() {
	return (
		<section className="min-h-screen bg-slate-950 px-6 py-10 text-slate-50">
			<div className="mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-6xl items-center justify-center">
				<div className="w-full max-w-md rounded-[2rem] border border-white/10 bg-white/[0.03] p-8 shadow-2xl shadow-black/40 backdrop-blur">
					<div className="space-y-3 text-center">
						<p className="text-xs font-medium uppercase tracking-[0.35em] text-cyan-200/80">
							Perfect Panel
						</p>
						<h1 className="text-2xl font-semibold tracking-tight text-white">正在恢复登录状态…</h1>
						<p className="text-sm leading-6 text-slate-300">
							请稍候，后台会先恢复你的会话，再继续进入目标页面。
						</p>
					</div>
				</div>
			</div>
		</section>
	);
}

export function RequireAuth({ children }: { children: ReactNode }) {
	const auth = useStore(authStore, (state) => state);
	const location = useRouterState({
		select: (state) => state.location,
	});
	const redirectTargetRef = useRef('/');

	if (location.pathname !== '/login') {
		redirectTargetRef.current = buildRedirectTarget(location);
	}

	useEffect(() => {
		if (!auth.hydrated) {
			hydrateAuthSession();
		}
	}, [auth.hydrated]);

	if (!auth.hydrated) {
		return <AuthHydrationFallback />;
	}

	if (!auth.accessToken) {
		return (
			<Navigate
				replace
				search={{
					redirect: redirectTargetRef.current,
				}}
				to="/login"
			/>
		);
	}

	return <>{children}</>;
}
