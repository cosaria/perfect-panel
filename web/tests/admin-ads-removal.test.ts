import { existsSync, readFileSync } from "node:fs";
import { describe, expect, test } from "bun:test";

function readRepoFile(relativePath: string) {
  return readFileSync(new URL(`../../${relativePath}`, import.meta.url), "utf8");
}

describe("admin ads removal", () => {
  test("removes ads page files and device show_ads toggle", () => {
    expect(
      existsSync(new URL("../apps/admin/src/pages/dashboard/ads/page.tsx", import.meta.url)),
    ).toBe(false);
    expect(
      existsSync(new URL("../apps/admin/src/pages/dashboard/ads/ads-form.tsx", import.meta.url)),
    ).toBe(false);
    expect(existsSync(new URL("../apps/admin/locales/en-US/ads.json", import.meta.url))).toBe(
      false,
    );
    expect(existsSync(new URL("../apps/admin/locales/zh-CN/ads.json", import.meta.url))).toBe(
      false,
    );

    const deviceForm = readRepoFile(
      "web/apps/admin/src/pages/dashboard/auth-control/forms/device-form.tsx",
    );
    expect(deviceForm).not.toContain("show_ads");
    expect(deviceForm).not.toContain("showAds");
  });
});
