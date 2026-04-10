import '@testing-library/jest-dom/vitest';

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { SETTINGS_DRAFT_STORAGE_KEY } from './settings-draft';
import { SettingsShell } from './settings-shell';

describe('SettingsShell', () => {
	beforeEach(() => {
		window.localStorage.clear();
	});

	afterEach(() => {
		cleanup();
	});

	it('loads an existing local draft into the form', async () => {
		window.localStorage.setItem(
			SETTINGS_DRAFT_STORAGE_KEY,
			JSON.stringify({
				siteName: 'Recovered Panel',
			}),
		);

		render(<SettingsShell />);

		await waitFor(() => {
			expect(screen.getByDisplayValue('Recovered Panel')).toBeInTheDocument();
		});
		expect(screen.getByText('已从本地恢复上次保存的设置草稿。')).toBeInTheDocument();
	});

	it('shows validation errors before saving an invalid draft', async () => {
		render(<SettingsShell />);

		fireEvent.change(screen.getByLabelText('站点名称'), {
			target: { value: '' },
		});
		fireEvent.change(screen.getByLabelText('站点地址'), {
			target: { value: 'not-a-url' },
		});
		fireEvent.click(screen.getByRole('button', { name: '保存草稿' }));

		await waitFor(() => {
			expect(screen.getByText('请输入站点名称。')).toBeInTheDocument();
		});
		expect(screen.getByText('请输入有效的站点地址。')).toBeInTheDocument();
	});

	it('persists a valid draft to localStorage and allows clearing it', async () => {
		render(<SettingsShell />);

		fireEvent.change(screen.getByLabelText('站点名称'), {
			target: { value: 'Panel Next' },
		});
		fireEvent.click(screen.getByRole('button', { name: '保存草稿' }));

		await waitFor(() => {
			expect(screen.getByText('本地草稿已保存，当前尚未同步到服务端。')).toBeInTheDocument();
		});

		expect(
			JSON.parse(window.localStorage.getItem(SETTINGS_DRAFT_STORAGE_KEY) ?? '{}'),
		).toMatchObject({
			siteName: 'Panel Next',
		});

		fireEvent.click(screen.getByRole('button', { name: '清空草稿' }));

		await waitFor(() => {
			expect(screen.getByText('本地草稿已清空，表单已恢复默认值。')).toBeInTheDocument();
		});
		expect(window.localStorage.getItem(SETTINGS_DRAFT_STORAGE_KEY)).toBeNull();
		expect(screen.getByLabelText('站点名称')).toHaveValue('Perfect Panel');
	});
});
