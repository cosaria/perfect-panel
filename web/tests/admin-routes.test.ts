import { describe, expect, test } from "bun:test";
import type { RouteObject } from "react-router-dom";

import { navs } from "../apps/admin/config/navs";
import { routes } from "../apps/admin/src/routes";

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

function collectNavUrls() {
  const urls = new Set<string>(["/", "/dashboard"]);

  for (const nav of navs) {
    if ("url" in nav) {
      urls.add(nav.url);
      continue;
    }

    for (const item of nav.items) {
      urls.add(item.url);
    }
  }

  return urls;
}

describe("admin routes", () => {
  test("covers every configured admin navigation path", () => {
    const actualPaths = collectRoutePaths(routes);
    const expectedPaths = collectNavUrls();

    for (const path of expectedPaths) {
      expect(actualPaths.has(path)).toBe(true);
    }
  });
});
