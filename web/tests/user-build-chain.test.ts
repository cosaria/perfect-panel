import { existsSync, readFileSync } from "node:fs";
import { describe, expect, test } from "bun:test";

import userPackageJson from "../apps/user/package.json";

function readRepoFile(relativePath: string) {
  return readFileSync(new URL(`../../${relativePath}`, import.meta.url), "utf8");
}

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

  test("embeds user assets from dist in the repo release chain", () => {
    const makefile = readRepoFile("Makefile");
    const dockerfile = readRepoFile("Dockerfile");
    const rootWebPackageJson = JSON.parse(readRepoFile("web/package.json")) as {
      scripts: Record<string, string>;
    };

    expect(makefile).toContain("cp -r web/apps/user/dist/* server/web/user-dist/");
    expect(dockerfile).toContain("COPY --from=web-builder /app/web/apps/user/dist ./web/user-dist");
    expect(dockerfile).toContain("COPY --from=builder /build/cache ./cache");
    expect(rootWebPackageJson.scripts.clean).toContain("apps/user/dist");
    expect(rootWebPackageJson.scripts.clean).not.toContain(".next");
    expect(rootWebPackageJson.scripts.clean).not.toContain("out");
  });

  test("drops user-only next build tooling files", () => {
    const userTsconfig = JSON.parse(readRepoFile("web/apps/user/tsconfig.json")) as {
      extends: string;
      include?: string[];
      compilerOptions?: { plugins?: Array<{ name?: string }> };
    };

    expect(userTsconfig.extends).toBe("@workspace/typescript-config/react-library.json");
    expect(userTsconfig.include ?? []).not.toContain("next-env.d.ts");
    expect(userTsconfig.include ?? []).not.toContain(".next/types/**/*.ts");
    expect(userTsconfig.compilerOptions?.plugins ?? []).toEqual([]);
    expect(existsSync(new URL("../apps/user/next.config.ts", import.meta.url))).toBe(false);
    expect(existsSync(new URL("../apps/user/next-env.d.ts", import.meta.url))).toBe(false);
  });
});
