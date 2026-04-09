import { existsSync } from "node:fs";
import { describe, expect, test } from "bun:test";
import type { RouteObject } from "react-router-dom";

import { navs } from "../apps/user/config/navs";

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

function findRoute(
  routeObjects: RouteObject[],
  matcher: (route: RouteObject) => boolean,
): RouteObject | null {
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
    const rendered = element.type(element.props as object) as {
      props?: Record<string, unknown>;
    } | null;
    return rendered?.props ?? null;
  }

  return (element as { props?: Record<string, unknown> }).props ?? null;
}

describe("user routes", () => {
  test("covers the unit 2 public and user navigation routes", async () => {
    const routesPath = new URL("../apps/user/src/routes.tsx", import.meta.url);

    expect(existsSync(routesPath)).toBe(true);

    const { routes } = (await import("../apps/user/src/routes")) as { routes: RouteObject[] };
    const actualPaths = collectRoutePaths(routes);

    const expectedPaths = new Set<string>([
      "/",
      "/auth",
      "/privacy-policy",
      "/tos",
      "/purchasing",
      "/purchasing/order",
      "/bind/:platform",
      "/oauth/:platform",
      "/dashboard",
      "/payment",
    ]);

    for (const nav of navs) {
      if ("url" in nav) {
        expectedPaths.add(nav.url);
        continue;
      }

      for (const item of nav.items) {
        expectedPaths.add(item.url);
      }
    }

    for (const path of expectedPaths) {
      expect(actualPaths.has(path)).toBe(true);
    }
  });

  test("renders an explicit 404 page for unknown public routes", async () => {
    const { routes } = (await import("../apps/user/src/routes")) as { routes: RouteObject[] };
    const shellRoute = routes[0]?.children?.find((route) => route.path === undefined) ?? null;
    const publicNotFoundRoute = shellRoute?.children?.find((route) => route.path === "*") ?? null;

    expect(publicNotFoundRoute).not.toBeNull();
    const props = resolveRouteElementProps(publicNotFoundRoute?.element);
    expect(props?.title).toBe("404");
    expect(props?.description).toContain("页面不存在");
  });

  test("renders an explicit 404 page for unknown dashboard routes", async () => {
    const { routes } = (await import("../apps/user/src/routes")) as { routes: RouteObject[] };
    const dashboardShellRoute = findRoute(
      routes,
      (route) =>
        route.element !== undefined && route.children?.some((child) => child.path === "dashboard"),
    );
    const dashboardNotFoundRoute =
      dashboardShellRoute?.children?.find((route) => route.path === "*") ?? null;

    expect(dashboardNotFoundRoute).not.toBeNull();
    const props = resolveRouteElementProps(dashboardNotFoundRoute?.element);
    expect(props?.title).toBe("404");
    expect(props?.description).toContain("页面不存在");
  });
});
