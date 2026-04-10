import '@testing-library/jest-dom/vitest';

import { render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it } from 'vitest';

import { Button } from './index';

describe('shared ui', () => {
	it('renders button text', () => {
		render(<Button>保存</Button>);

		expect(screen.getByText('保存')).toBeInTheDocument();
	});
});
