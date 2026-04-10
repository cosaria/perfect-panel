import '@testing-library/jest-dom/vitest';

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { AppProviders } from './providers';

describe('admin providers', () => {
	it('renders children inside providers', () => {
		render(
			<AppProviders>
				<div>admin-shell</div>
			</AppProviders>,
		);

		expect(screen.getByText('admin-shell')).toBeInTheDocument();
	});
});
