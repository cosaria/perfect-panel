import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..', '..');

describe('build chain scripts', () => {
	it('defines root scripts and includes all smoke tests', () => {
		const pkg = JSON.parse(readFileSync(resolve(root, 'package.json'), 'utf8')) as {
			scripts: Record<string, string>;
		};

		expect(pkg.scripts.openapi).toBeDefined();
		expect(pkg.scripts.build).toBeDefined();
		expect(pkg.scripts.typecheck).toBeDefined();
		expect(pkg.scripts.test).toContain('tests/smoke/*.test.ts');
	});
});
