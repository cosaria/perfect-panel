import { describe, expect, test } from "bun:test";
import { buildApiBaseUrl, getApiRoot } from "../packages/ui/src/utils/api-base";

describe("api-base", () => {
  test("getApiRoot prefers explicit api url", () => {
    expect(getApiRoot("http://api.example.com", "http://site.example.com")).toBe(
      "http://api.example.com",
    );
  });

  test("buildApiBaseUrl keeps the api prefix for common routes", () => {
    expect(buildApiBaseUrl("", "", "/api/v1/common")).toBe("/api/v1/common");
  });

  test("buildApiBaseUrl trims a trailing slash from the root", () => {
    expect(buildApiBaseUrl("http://localhost:8080/", "", "/api/v1/admin")).toBe(
      "http://localhost:8080/api/v1/admin",
    );
  });
});
