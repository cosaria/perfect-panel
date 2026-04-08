import { existsSync } from "node:fs";
import { describe, expect, test } from "bun:test";

import userPackageJson from "../apps/user/package.json";

describe("user build chain", () => {
  test("uses vite-native user scripts and dist output", async () => {
    const viteConfigPath = new URL("../apps/user/vite.config.ts", import.meta.url);

    expect(existsSync(viteConfigPath)).toBe(true);
    expect(userPackageJson.scripts.dev).toBe("vite");
    expect(userPackageJson.scripts.build).toBe("vite build");
    expect(userPackageJson.scripts.preview).toContain("vite preview");

    const userViteConfig = (await import("../apps/user/vite.config")).default;
    expect(userViteConfig.build?.outDir).toBe("dist");
  });
});
