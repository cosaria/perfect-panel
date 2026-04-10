import { Outlet, createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '../shared/auth/require-auth';
import { AdminSidebar } from '../widgets/admin-sidebar';
import { AdminTopbar } from '../widgets/admin-topbar';

export const Route = createFileRoute('/_authed')({
	component: AuthedLayout,
});

function AuthedLayout() {
	return (
		<RequireAuth>
			<div className="min-h-screen bg-slate-950 text-slate-50 lg:grid lg:grid-cols-[18rem_minmax(0,1fr)]">
				<AdminSidebar />
				<div className="min-w-0">
					<AdminTopbar />
					<main className="p-6 lg:p-8">
						<Outlet />
					</main>
				</div>
			</div>
		</RequireAuth>
	);
}
