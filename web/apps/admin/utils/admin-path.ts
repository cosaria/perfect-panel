import { NEXT_PUBLIC_ADMIN_PATH } from "@/config/constants";

export const ADMIN_HOME_PATH = "/dashboard/workplace";

function hasPathPrefix(pathname: string, prefix: string) {
  return pathname === prefix || pathname.startsWith(`${prefix}/`);
}

export function normalizeAdminPath(value?: string) {
  const trimmed = value?.trim() ?? "";
  if (!trimmed || trimmed === "/") {
    return "/admin";
  }

  const normalized = `/${trimmed.replace(/^\/+/, "").replace(/\/+$/, "")}`;
  return normalized === "/" ? "/admin" : normalized;
}

export function stripAdminPath(pathname: string, adminPath: string = NEXT_PUBLIC_ADMIN_PATH) {
  const normalizedAdminPath = normalizeAdminPath(adminPath);
  if (!hasPathPrefix(pathname, normalizedAdminPath)) {
    return pathname || "/";
  }

  const stripped = pathname.slice(normalizedAdminPath.length);
  return stripped || "/";
}

export function toAdminPath(pathname: string, adminPath: string = NEXT_PUBLIC_ADMIN_PATH) {
  const normalizedAdminPath = normalizeAdminPath(adminPath);
  const normalizedPathname = pathname.trim() || "/";

  if (/^https?:\/\//.test(normalizedPathname)) {
    return normalizedPathname;
  }

  if (
    normalizedPathname === normalizedAdminPath ||
    normalizedPathname.startsWith(`${normalizedAdminPath}/`)
  ) {
    return normalizedPathname;
  }

  if (normalizedPathname === "/admin" || normalizedPathname.startsWith("/admin/")) {
    return (
      `${normalizedAdminPath}${normalizedPathname.slice("/admin".length)}` || normalizedAdminPath
    );
  }

  if (normalizedPathname === "/" || normalizedPathname === "") {
    return normalizedAdminPath;
  }

  const withLeadingSlash = normalizedPathname.startsWith("/")
    ? normalizedPathname
    : `/${normalizedPathname}`;

  return `${normalizedAdminPath}${withLeadingSlash}`;
}

export function canonicalizeAdminBrowserPath(
  pathname: string,
  adminPath: string = NEXT_PUBLIC_ADMIN_PATH,
) {
  if (pathname === "/admin" || pathname.startsWith("/admin/")) {
    return toAdminPath(pathname, adminPath);
  }

  return pathname;
}
