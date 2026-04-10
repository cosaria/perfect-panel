import { existsSync, readdirSync } from 'node:fs';
import { resolve } from 'node:path';
import process from 'node:process';
import { createClient, defaultPlugins } from '@hey-api/openapi-ts';

const packageRoot = resolve(import.meta.dirname, '..');
const repoRoot = resolve(packageRoot, '..', '..', '..');
const output = resolve(packageRoot, 'src', 'generated');

const inputCandidates = buildInputCandidates({
	repoRoot,
});
const input = inputCandidates.find((candidate) => existsSync(candidate));

if (!input) {
	console.error(
		[
			'缺少 OpenAPI 输入文件，已尝试以下路径：',
			...inputCandidates.map((candidate) => `- ${candidate}`),
			'可以通过环境变量 SERVER_V2_OPENAPI_INPUT 显式指定输入文件。',
		].join('\n'),
	);
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

function buildInputCandidates({ repoRoot }) {
	const candidates = [];
	const envInput = process.env.SERVER_V2_OPENAPI_INPUT?.trim();

	if (envInput) {
		candidates.push(resolve(envInput));
	}

	candidates.push(
		resolve(repoRoot, '..', 'perfect-panel', 'server-v2', 'openapi', 'dist', 'openapi.json'),
	);

	const workspaceRoot = resolve(repoRoot, '..');

	for (const entry of readdirSync(workspaceRoot, { withFileTypes: true })) {
		if (!entry.isDirectory()) continue;
		if (entry.name === 'perfect-panel-web-v2') continue;

		const candidate = resolve(
			workspaceRoot,
			entry.name,
			'server-v2',
			'openapi',
			'dist',
			'openapi.json',
		);

		candidates.push(candidate);
	}

	return [...new Set(candidates)];
}
