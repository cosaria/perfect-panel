import { existsSync } from "node:fs";
import { describe, expect, test } from "bun:test";
import type { RouteObject } from "react-router-dom";

function joinRoutePath(parentPath: string, childPath?: string) {
  if (!childPath) return parentPath || "/";
  if (childPath.startsWith("/")) return childPath;
  if (!parentPath || parentPath === "/") return `/${childPath}`;
  return `${parentPath}/${childPath}`;
}

function collectRoutePaths(routeObjects: RouteObject[], parentPath = ""): Set<string> {
  const paths = new Set<string>();

  for (const route of routeObjects) {
    if (route.index) {
      paths.add(parentPath || "/");
      continue;
    }

    if (route.path === "*") {
      continue;
    }

    const currentPath = joinRoutePath(parentPath, route.path);
    if (route.path) {
      paths.add(currentPath);
    }

    if (route.children) {
      const childPaths = collectRoutePaths(route.children, currentPath);
      for (const childPath of childPaths) {
        paths.add(childPath);
      }
    }
  }

  return paths;
}

describe("user routes", () => {
  test("covers the unit 1 spike routes", async () => {
    const routesPath = new URL("../apps/user/src/routes.tsx", import.meta.url);

    expect(existsSync(routesPath)).toBe(true);

    const { routes } = (await import("../apps/user/src/routes")) as { routes: RouteObject[] };
    const actualPaths = collectRoutePaths(routes);

    expect(actualPaths.has("/")).toBe(true);
    expect(actualPaths.has("/auth")).toBe(true);
    expect(actualPaths.has("/dashboard")).toBe(true);
  });
});
