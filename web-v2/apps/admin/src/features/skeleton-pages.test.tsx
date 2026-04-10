import '@testing-library/jest-dom/vitest';

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { SettingsShell } from './settings/settings-shell';

describe('skeleton pages', () => {
	it('shows the second-stage settings draft shell', () => {
		render(<SettingsShell />);

		expect(screen.getByText('系统设置')).toBeInTheDocument();
		expect(
			screen.getByText('当前为前端草稿模式，保存后只会写入本地浏览器，不会提交到服务端。'),
		).toBeInTheDocument();
	});
});
