import { existsSync } from 'node:fs';
import { resolve } from 'node:path';
import process from 'node:process';
import { createClient, defaultPlugins } from '@hey-api/openapi-ts';

const packageRoot = resolve(import.meta.dirname, '..');
const repoRoot = resolve(packageRoot, '..', '..', '..');
const input = resolve(
	repoRoot,
	'..',
	'..',
	'perfect-panel',
	'server-v2',
	'openapi',
	'dist',
	'openapi.json',
);
const output = resolve(packageRoot, 'src', 'generated');

if (!existsSync(input)) {
	console.error(`缺少 OpenAPI 输入文件: ${input}`);
	process.exit(1);
}

await createClient({
	input,
	output: {
		path: output,
		clean: true,
	},
	plugins: [
		...defaultPlugins,
		'@hey-api/client-fetch',
		{
			name: '@hey-api/typescript',
			enums: 'javascript',
		},
		{
			name: '@hey-api/sdk',
			auth: false,
		},
	],
});
