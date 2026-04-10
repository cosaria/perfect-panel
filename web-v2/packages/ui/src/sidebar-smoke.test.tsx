import '@testing-library/jest-dom/vitest';

import { render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it } from 'vitest';

import { Button, Sidebar } from './index';

describe('shared ui', () => {
	it('renders button text', () => {
		render(<Button>保存</Button>);

		expect(screen.getByText('保存')).toBeInTheDocument();
	});

	it('renders sidebar export', () => {
		render(<Sidebar aria-label="主侧边栏">导航</Sidebar>);

		expect(screen.getByRole('complementary', { name: '主侧边栏' })).toBeInTheDocument();
		expect(screen.getByText('导航')).toBeInTheDocument();
	});
});
