import { readFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..', '..');
const expectedScripts = {
	dev: 'turbo dev',
	build: 'turbo build',
	test: 'vitest run tests/smoke/workspace-layout.test.ts && turbo test',
	lint: 'turbo lint',
	format: 'biome format . --write',
	typecheck: 'turbo typecheck',
	openapi: 'turbo openapi',
};

const expectedTurboTasks = {
	dev: {
		cache: false,
		persistent: true,
	},
	build: {
		dependsOn: ['^build'],
		outputs: ['dist/**', '.tanstack/**'],
	},
	test: {
		dependsOn: ['^test'],
		outputs: [],
	},
	lint: {
		dependsOn: ['^lint'],
		outputs: [],
	},
	typecheck: {
		dependsOn: ['^typecheck'],
		outputs: [],
	},
	openapi: {
		dependsOn: ['^openapi'],
		outputs: ['src/generated/**'],
	},
};

describe('workspace layout', () => {
	it('matches the task 1 monorepo contract', () => {
		const packageJson = JSON.parse(readFileSync(resolve(root, 'package.json'), 'utf8')) as {
			private: boolean;
			packageManager: string;
			scripts: Record<string, string>;
		};
		const turboJson = JSON.parse(readFileSync(resolve(root, 'turbo.json'), 'utf8')) as {
			tasks: Record<string, unknown>;
		};
		const workspaceLines = readFileSync(resolve(root, 'pnpm-workspace.yaml'), 'utf8')
			.trim()
			.split('\n')
			.map((line) => line.trim());

		expect(packageJson.private).toBe(true);
		expect(packageJson.packageManager).toBe('pnpm@10.6.0');
		expect(packageJson.scripts).toMatchObject(expectedScripts);
		expect(turboJson.tasks).toEqual(expectedTurboTasks);
		expect(workspaceLines).toEqual(['packages:', "- 'apps/*'", "- 'packages/*'"]);
	});
});
