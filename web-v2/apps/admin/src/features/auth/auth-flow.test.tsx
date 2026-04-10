import '@testing-library/jest-dom/vitest';

import { fireEvent, render, screen } from '@testing-library/react';
import { RouterContextProvider } from '@tanstack/react-router';
import { describe, expect, it } from 'vitest';
import { createRouter } from '../../router';
import { LoginForm } from './login-form';

describe('auth flow', () => {
	it('renders login fields and allows typing', () => {
		const router = createRouter();

		render(
			<RouterContextProvider router={router}>
				<LoginForm />
			</RouterContextProvider>,
		);

		expect(screen.getByText('邮箱')).toBeInTheDocument();
		expect(screen.getByText('密码')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument();

		fireEvent.change(screen.getByLabelText('邮箱'), { target: { value: 'admin@example.com' } });
		fireEvent.change(screen.getByLabelText('密码'), { target: { value: 'secret123' } });

		expect(screen.getByLabelText('邮箱')).toHaveValue('admin@example.com');
		expect(screen.getByLabelText('密码')).toHaveValue('secret123');
	});
});
