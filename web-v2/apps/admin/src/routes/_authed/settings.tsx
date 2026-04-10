import { createFileRoute } from '@tanstack/react-router';
import { SettingsShell } from '../../features/settings/settings-shell';

export const Route = createFileRoute('/_authed/settings')({
	component: SettingsRoute,
});

function SettingsRoute() {
	return <SettingsShell />;
}
