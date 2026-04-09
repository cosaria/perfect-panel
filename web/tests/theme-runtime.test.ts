import { describe, expect, test } from "bun:test";
import {
  applyResolvedTheme,
  readStoredTheme,
  resolveSonnerTheme,
  resolveTheme,
  type Theme,
} from "../packages/ui/src/lib/theme";

function createClassList(initialClasses: string[] = []) {
  const classes = new Set(initialClasses);

  return {
    add(...tokens: string[]) {
      for (const token of tokens) {
        classes.add(token);
      }
    },
    contains(token: string) {
      return classes.has(token);
    },
    remove(...tokens: string[]) {
      for (const token of tokens) {
        classes.delete(token);
      }
    },
    toArray() {
      return [...classes].sort();
    },
  };
}

function createStorage(initialValue?: string) {
  const values = new Map<string, string>();

  if (initialValue) {
    values.set("theme", initialValue);
  }

  return {
    getItem(key: string) {
      return values.get(key) ?? null;
    },
    setItem(key: string, value: string) {
      values.set(key, value);
    },
  };
}

describe("theme runtime", () => {
  test("resolveTheme uses the system theme when the selected theme is system", () => {
    expect(resolveTheme("system", "dark")).toBe("dark");
    expect(resolveTheme("light", "dark")).toBe("light");
  });

  test("applyResolvedTheme keeps only the active color class and color-scheme", () => {
    const classList = createClassList(["app", "light"]);
    const element = {
      classList,
      style: {
        colorScheme: "light",
      },
    };

    applyResolvedTheme(element, "dark");

    expect(classList.toArray()).toEqual(["app", "dark"]);
    expect(element.style.colorScheme).toBe("dark");
  });

  test("resolveSonnerTheme maps system theme to the resolved theme", () => {
    expect(resolveSonnerTheme("system", "dark")).toBe("dark");
    expect(resolveSonnerTheme("light", "dark")).toBe("light");
  });

  test("readStoredTheme falls back to the provided default when storage is invalid", () => {
    const storage = createStorage("sepia");

    expect(readStoredTheme(storage, "theme", "system")).toBe("system");
  });

  test("readStoredTheme returns the stored theme when it is supported", () => {
    const storage = createStorage("dark");

    expect(readStoredTheme(storage, "theme", "system")).toBe("dark" satisfies Theme);
  });
});
