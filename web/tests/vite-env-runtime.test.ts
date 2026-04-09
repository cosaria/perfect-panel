import { describe, expect, test } from "bun:test";

function runBunEval(script: string) {
  const result = Bun.spawnSync({
    cmd: ["bun", "--eval", script],
    cwd: new URL("..", import.meta.url).pathname,
    stderr: "pipe",
    stdout: "pipe",
  });

  return {
    exitCode: result.exitCode,
    stderr: new TextDecoder().decode(result.stderr).trim(),
    stdout: new TextDecoder().decode(result.stdout).trim(),
  };
}

describe("vite env runtime", () => {
  test("admin redirect logic reads VITE_ADMIN_PATH from window.__ENV", () => {
    const result = runBunEval(`
      globalThis.window = { __ENV: { VITE_ADMIN_PATH: "/manage" } };
      globalThis.sessionStorage = {
        getItem() { return "/dashboard"; },
        removeItem() {},
        setItem() {},
      };
      const { getRedirectUrl } = await import("./apps/admin/utils/common.ts");
      console.log(getRedirectUrl());
    `);

    expect(result.exitCode).toBe(0);
    expect(result.stderr).toBe("");
    expect(result.stdout).toBe("/manage/dashboard");
  });

  test("user setup-clients prefers VITE_API_URL when configuring common-api", () => {
    const result = runBunEval(`
      globalThis.window = {
        __ENV: {
          VITE_API_URL: "http://api.example.com",
          VITE_SITE_URL: "http://site.example.com",
        },
      };
      await import("./apps/user/utils/setup-clients.ts");
      const { client: userCommonClient } = await import("./apps/user/services/common-api/client.gen.ts");
      console.log(userCommonClient.getConfig().baseUrl);
    `);

    expect(result.exitCode).toBe(0);
    expect(result.stderr).toBe("");
    expect(result.stdout).toBe("http://api.example.com/api/v1/common");
  });
});
