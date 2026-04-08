import { describe, expect, test } from "bun:test";

import {
  canonicalizeAdminBrowserPath,
  normalizeAdminPath,
  stripAdminPath,
  toAdminPath,
} from "../apps/admin/utils/admin-path";

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

  test("only strips the configured admin prefix on a real path boundary", () => {
    expect(stripAdminPath("/manage/dashboard", "/manage")).toBe("/dashboard");
    expect(stripAdminPath("/management/dashboard", "/manage")).toBe("/management/dashboard");
  });
});
