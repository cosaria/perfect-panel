import { describe, expect, test } from "bun:test";

const NEXT_IMPORT_PATTERN =
  'from ["\']next/(navigation|link|image|legacy/image)["\']|from ["\']@/src/compat/(app|next)-(navigation|link|image)["\']';

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
});
