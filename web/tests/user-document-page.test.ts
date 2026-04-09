import { existsSync, readFileSync } from "node:fs";
import { describe, expect, test } from "bun:test";

function readRepoFile(relativePath: string) {
  return readFileSync(new URL(`../../${relativePath}`, import.meta.url), "utf8");
}

describe("user document page", () => {
  test("renders only admin-managed documents and does not depend on bundled tutorials", () => {
    const pageSource = readRepoFile("web/apps/user/src/pages/(main)/(user)/document/page.tsx");

    expect(pageSource).toContain("queryDocumentList");
    expect(pageSource).toContain("DocumentButton");

    expect(pageSource).not.toContain("getTutorialList");
    expect(pageSource).not.toContain("TutorialButton");
    expect(pageSource).not.toContain("NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT");

    expect(
      existsSync(
        new URL(
          "../apps/user/src/pages/(main)/(user)/document/tutorial-button.tsx",
          import.meta.url,
        ),
      ),
    ).toBe(false);
    expect(existsSync(new URL("../apps/user/utils/tutorial.ts", import.meta.url))).toBe(false);
  });
});
