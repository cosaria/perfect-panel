import { afterEach, beforeEach, describe, expect, test } from "bun:test";

const USER_LOCALE_CHANGE_EVENT = "ppanel:user-locale-change";

const originalDocument = globalThis.document;
const originalNavigator = globalThis.navigator;
const originalWindow = globalThis.window;

function createWindowLike() {
  const eventTarget = new EventTarget();
  const storage = new Map<string, string>();

  return Object.assign(eventTarget, {
    localStorage: {
      getItem(key: string) {
        return storage.get(key) ?? null;
      },
      removeItem(key: string) {
        storage.delete(key);
      },
      setItem(key: string, value: string) {
        storage.set(key, value);
      },
    },
    location: {
      reload() {},
    },
  });
}

describe("user locale runtime", () => {
  beforeEach(() => {
    Object.defineProperty(globalThis, "document", {
      configurable: true,
      value: { cookie: "" },
      writable: true,
    });
    Object.defineProperty(globalThis, "navigator", {
      configurable: true,
      value: { language: "en-US" },
      writable: true,
    });
    Object.defineProperty(globalThis, "window", {
      configurable: true,
      value: createWindowLike(),
      writable: true,
    });
  });

  afterEach(() => {
    Object.defineProperty(globalThis, "document", {
      configurable: true,
      value: originalDocument,
      writable: true,
    });
    Object.defineProperty(globalThis, "navigator", {
      configurable: true,
      value: originalNavigator,
      writable: true,
    });
    Object.defineProperty(globalThis, "window", {
      configurable: true,
      value: originalWindow,
      writable: true,
    });
  });

  test("setLocale updates the cookie-backed locale and publishes a runtime change event", async () => {
    const { getClientLocale } = await import("../apps/user/locales/client");
    const { setLocale } = await import("../apps/user/utils/common");
    const receivedLocales: string[] = [];

    window.addEventListener(USER_LOCALE_CHANGE_EVENT, (event) => {
      receivedLocales.push((event as CustomEvent<{ locale: string }>).detail.locale);
    });

    setLocale("zh-CN");

    expect(getClientLocale()).toBe("zh-CN");
    expect(receivedLocales).toEqual(["zh-CN"]);
  });
});
