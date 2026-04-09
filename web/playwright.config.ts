import { defineConfig } from "@playwright/test";

const port = Number(process.env.STATIC_SMOKE_PORT ?? "4173");
const adminPath = process.env.STATIC_SMOKE_ADMIN_PATH ?? "/manage";

export default defineConfig({
  testDir: "./tests/smoke",
  fullyParallel: false,
  retries: process.env.CI ? 2 : 0,
  reporter: "list",
  use: {
    baseURL: `http://127.0.0.1:${port}`,
    trace: "retain-on-failure",
  },
  webServer: {
    command: `cd ../server && STATIC_SMOKE_ADDR=127.0.0.1:${port} STATIC_SMOKE_SITE_URL=http://127.0.0.1:${port} STATIC_SMOKE_ADMIN_PATH=${adminPath} go run -tags embed ./tests/static_smoke`,
    port,
    reuseExistingServer: !process.env.CI,
    stdout: "pipe",
    stderr: "pipe",
    timeout: 120_000,
  },
});
