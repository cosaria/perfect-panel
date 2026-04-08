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
  const urls = new Set<string>(["/", "/dashboard/workplace"]);

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

function findRoute(routeObjects: RouteObject[], matcher: (route: RouteObject) => boolean): RouteObject | null {
  for (const route of routeObjects) {
    if (matcher(route)) {
      return route;
    }

    if (route.children) {
      const matched = findRoute(route.children, matcher);
      if (matched) {
        return matched;
      }
    }
  }

  return null;
}

function resolveRouteElementProps(element: RouteObject["element"]) {
  if (!element || typeof element !== "object" || !("type" in element) || !("props" in element)) {
    return null;
  }

  if (typeof element.type === "function") {
    const rendered = element.type(element.props as object) as { props?: Record<string, unknown> } | null;
    return rendered?.props ?? null;
  }

  return (element as { props?: Record<string, unknown> }).props ?? null;
}

describe("admin routes", () => {
  test("covers every configured admin navigation path", () => {
    const actualPaths = collectRoutePaths(routes);
    const expectedPaths = collectNavUrls();

    expect(actualPaths.has("/dashboard/workplace")).toBe(true);
    expect(expectedPaths.has("/dashboard/workplace")).toBe(true);

    expect(actualPaths.has("/dashboard/ads")).toBe(false);
    expect(expectedPaths.has("/dashboard/ads")).toBe(false);

    for (const path of expectedPaths) {
      expect(actualPaths.has(path)).toBe(true);
    }
  });

  test("renders an explicit 404 page for the legacy bare dashboard route", () => {
    const dashboardRoute = findRoute(routes, (route) => route.path === "dashboard");
    const dashboardIndexRoute = dashboardRoute?.children?.find((route) => route.index) ?? null;

    expect(dashboardIndexRoute).not.toBeNull();

    const props = resolveRouteElementProps(dashboardIndexRoute?.element);
    expect(props?.title).toBe("404");
    expect(props?.description).toContain("页面不存在");
  });

  test("renders an explicit 404 page for unknown dashboard routes", () => {
    const dashboardRoute = findRoute(routes, (route) => route.path === "dashboard");
    const dashboardNotFoundRoute = dashboardRoute?.children?.find((route) => route.path === "*") ?? null;

    expect(dashboardNotFoundRoute).not.toBeNull();
    const props = resolveRouteElementProps(dashboardNotFoundRoute?.element);
    expect(props?.title).toBe("404");
    expect(props?.description).toContain("页面不存在");
  });

  test("renders an explicit 404 page for unknown top-level admin routes", () => {
    const rootRoute = routes[0];
    const topLevelNotFoundRoute = rootRoute.children?.find((route) => route.path === "*") ?? null;

    expect(topLevelNotFoundRoute).not.toBeNull();
    const props = resolveRouteElementProps(topLevelNotFoundRoute?.element);
    expect(props?.title).toBe("404");
    expect(props?.description).toContain("页面不存在");
  });
});
