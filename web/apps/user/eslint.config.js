import { nextJsConfig } from "@workspace/eslint-config/next-js";

/** @type {import("eslint").Linter.Config} */
export default [
	{
		ignores: ["next-env.d.ts", "services/**/*.ts"],
	},
	...nextJsConfig,
	{
		rules: {
			"@typescript-eslint/no-namespace": "off",
			"@typescript-eslint/no-explicit-any": "off",
			"@typescript-eslint/no-unused-vars": "off",
		},
	},
];
