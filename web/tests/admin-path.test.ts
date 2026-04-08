import { describe, expect, test } from "bun:test";

import {
  canonicalizeAdminBrowserPath,
  normalizeAdminPath,
  stripAdminPath,
  toAdminPath,
} from "../apps/admin/utils/admin-path";
import { getRedirectUrl } from "../apps/admin/utils/common";

describe("admin-path", () => {
  test("normalizes the configured admin base path", () => {
    expect(normalizeAdminPath()).toBe("/admin");
    expect(normalizeAdminPath("manage")).toBe("/manage");
    expect(normalizeAdminPath("/manage/")).toBe("/manage");
  });

  test("prefixes dashboard routes with the runtime admin path", () => {
    expect(toAdminPath("/dashboard", "/manage")).toBe("/manage/dashboard");
    expect(toAdminPath("/manage/dashboard", "/manage")).toBe("/manage/dashboard");
    expect(toAdminPath("/dashboard/order?search=123", "/manage")).toBe(
      "/manage/dashboard/order?search=123",
    );
  });

  test("maps the admin root to the configured entry path", () => {
    expect(toAdminPath("/", "/manage")).toBe("/manage");
    expect(toAdminPath("", "/manage")).toBe("/manage");
  });

  test("canonicalizes legacy admin browser paths to the runtime admin path", () => {
    expect(canonicalizeAdminBrowserPath("/admin/dashboard", "/manage")).toBe("/manage/dashboard");
    expect(canonicalizeAdminBrowserPath("/manage/dashboard", "/manage")).toBe("/manage/dashboard");
  });

  test("defaults auth redirect targets to the workplace home route", () => {
    const originalSessionStorage = globalThis.sessionStorage;

    Object.defineProperty(globalThis, "sessionStorage", {
      configurable: true,
      value: {
        getItem() {
          return null;
        },
      },
      writable: true,
    });

    expect(getRedirectUrl()).toBe("/admin/dashboard/workplace");

    Object.defineProperty(globalThis, "sessionStorage", {
      configurable: true,
      value: originalSessionStorage,
      writable: true,
    });
  });

  test("keeps legacy dashboard redirect targets unchanged", () => {
    const originalSessionStorage = globalThis.sessionStorage;

    Object.defineProperty(globalThis, "sessionStorage", {
      configurable: true,
      value: {
        getItem() {
          return "/dashboard";
        },
      },
      writable: true,
    });

    expect(getRedirectUrl()).toBe("/admin/dashboard");

    Object.defineProperty(globalThis, "sessionStorage", {
      configurable: true,
      value: originalSessionStorage,
      writable: true,
    });
  });

  test("only strips the configured admin prefix on a real path boundary", () => {
    expect(stripAdminPath("/manage/dashboard", "/manage")).toBe("/dashboard");
    expect(stripAdminPath("/management/dashboard", "/manage")).toBe("/management/dashboard");
  });
});
