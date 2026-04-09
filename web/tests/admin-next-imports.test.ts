import { describe, expect, test } from "bun:test";
import { existsSync } from "node:fs";

const NEXT_IMPORT_PATTERN =
	"from [\"']next/(navigation|link|image|legacy/image)[\"']|from [\"']@/src/compat/.+[\"']|[\"']@/app/";

describe("admin next import cleanup", () => {
	test("admin source no longer imports next navigation, link, or image modules", () => {
		const result = Bun.spawnSync({
			cmd: ["rg", "-n", NEXT_IMPORT_PATTERN, "web/apps/admin"],
			cwd: "/Users/admin/.config/superpowers/worktrees/perfect-panel/feat-web-vite-only-migration",
			stderr: "pipe",
			stdout: "pipe",
		});

		const output = result.stdout.toString().trim();
		expect(output).toBe("");
	});

	test("admin app directory has been retired in favor of src/pages", () => {
		expect(existsSync(new URL("../apps/admin/app", import.meta.url))).toBe(
			false,
		);
		expect(
			existsSync(new URL("../apps/admin/src/pages", import.meta.url)),
		).toBe(true);
	});

	test("admin source no longer keeps a src/compat directory", () => {
		expect(
			existsSync(new URL("../apps/admin/src/compat", import.meta.url)),
		).toBe(false);
	});
});
