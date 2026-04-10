import '@testing-library/jest-dom/vitest';

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { UsersTableShell } from './users/users-table-shell';

describe('skeleton pages', () => {
	it('shows empty-state structure for users page', () => {
		render(<UsersTableShell />);

		expect(screen.getByText('用户管理')).toBeInTheDocument();
		expect(screen.getByText('当前阶段先提供页面结构和空状态。')).toBeInTheDocument();
	});
});
