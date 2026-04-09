import { existsSync } from "node:fs";
import { describe, expect, test } from "bun:test";

const NEXT_IMPORT_PATTERN =
  "from [\"']next/(navigation|link|image|legacy/image)[\"']|from [\"']@/src/compat/.+[\"']|[\"']@/app/";

describe("user next import cleanup", () => {
  test("user source no longer imports next navigation, link, or image modules", () => {
    const result = Bun.spawnSync({
      cmd: ["rg", "-n", NEXT_IMPORT_PATTERN, "web/apps/user"],
      cwd: "/Users/admin/.config/superpowers/worktrees/perfect-panel/feat-web-vite-only-migration",
      stderr: "pipe",
      stdout: "pipe",
    });

    const output = result.stdout.toString().trim();
    expect(output).toBe("");
  });

  test("user app directory has been retired in favor of src/pages", () => {
    expect(existsSync(new URL("../apps/user/app", import.meta.url))).toBe(false);
    expect(existsSync(new URL("../apps/user/src/pages", import.meta.url))).toBe(true);
  });

  test("user source no longer depends on app compat wrappers", () => {
    expect(existsSync(new URL("../apps/user/src/compat", import.meta.url))).toBe(false);
  });
});
