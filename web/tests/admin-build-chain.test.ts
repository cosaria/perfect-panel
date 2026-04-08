import { existsSync, readFileSync } from "node:fs";
import { describe, expect, test } from "bun:test";

import adminPackageJson from "../apps/admin/package.json";
import turboJson from "../turbo.json";
import viteConfig from "../apps/admin/vite.config";

function readRepoFile(relativePath: string) {
  return readFileSync(new URL(`../../${relativePath}`, import.meta.url), "utf8");
}

describe("admin build chain", () => {
  test("uses vite-native admin scripts and dist output", () => {
    expect(adminPackageJson.scripts.dev).toBe("vite");
    expect(adminPackageJson.scripts.build).toBe("vite build");
    expect(adminPackageJson.scripts.preview).toContain("vite preview");
    expect(viteConfig.build?.outDir).toBe("dist");
  });

  test("embeds admin assets from dist in the repo release chain", () => {
    const makefile = readRepoFile("Makefile");
    const dockerfile = readRepoFile("Dockerfile");
    const rootWebPackageJson = JSON.parse(readRepoFile("web/package.json")) as {
      scripts: Record<string, string>;
    };

    expect(makefile).toContain("cp -r web/apps/admin/dist/* server/web/admin-dist/");
    expect(dockerfile).toContain(
      "COPY --from=web-builder /app/web/apps/admin/dist ./web/admin-dist",
    );
    expect(rootWebPackageJson.scripts.clean).toContain("apps/admin/dist");
  });

  test("locks frontend dependencies during docker builds", () => {
    const dockerfile = readRepoFile("Dockerfile");

    expect(dockerfile).toContain("RUN bun install --frozen-lockfile");
  });

  test("ignores vite dist assets in workspace lint inputs", () => {
    const webGitignore = readRepoFile("web/.gitignore");

    expect(webGitignore).toContain("apps/*/dist");
  });

  test("declares both vite and next build outputs to turbo", () => {
    expect(turboJson.tasks.build.outputs).toContain("dist/**");
    expect(turboJson.tasks.build.outputs).toContain("out/**");
    expect(turboJson.tasks.build.outputs).toContain(".next/**");
  });

  test("drops admin-only next build tooling files", () => {
    const adminTsconfig = JSON.parse(readRepoFile("web/apps/admin/tsconfig.json")) as {
      extends: string;
      include?: string[];
      compilerOptions?: { plugins?: Array<{ name?: string }> };
    };

    expect(adminTsconfig.extends).toBe("@workspace/typescript-config/react-library.json");
    expect(adminTsconfig.include ?? []).not.toContain("next-env.d.ts");
    expect(adminTsconfig.include ?? []).not.toContain(".next/types/**/*.ts");
    expect(adminTsconfig.compilerOptions?.plugins ?? []).toEqual([]);
    expect(existsSync(new URL("../apps/admin/next.config.ts", import.meta.url))).toBe(false);
    expect(existsSync(new URL("../apps/admin/next-env.d.ts", import.meta.url))).toBe(false);
  });
});
