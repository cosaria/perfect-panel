import '@testing-library/jest-dom/vitest';

import { RouterContextProvider } from '@tanstack/react-router';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { createRouter } from '../router';
import { AdminSidebar } from './admin-sidebar';

describe('AdminSidebar', () => {
	it('renders the first-stage navigation items', () => {
		const router = createRouter();

		render(
			<RouterContextProvider router={router}>
				<AdminSidebar />
			</RouterContextProvider>,
		);

		expect(screen.getByText('Dashboard')).toBeInTheDocument();
		expect(screen.getByText('用户管理')).toBeInTheDocument();
		expect(screen.getByText('系统设置')).toBeInTheDocument();
	});
});
