import { describe, expect, test } from "bun:test";
import { createTranslator } from "../packages/ui/src/components/i18n-provider";

describe("i18n translator", () => {
  test("createTranslator translates keys from a namespace message object", () => {
    const t = createTranslator({
      locale: "en-US",
      messages: {
        submit: "Submit",
        title: "Sign in",
      },
    });

    expect(t("title")).toBe("Sign in");
    expect(t("submit")).toBe("Submit");
  });
});
