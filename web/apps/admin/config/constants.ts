import packageJSON from "../package.json";

declare global {
  interface Window {
    __ENV?: Record<string, string>;
  }
}

function getEnv(key: string): string | undefined {
  if (typeof window !== "undefined") {
    return window.__ENV?.[key];
  }
  return undefined;
}

export const locales = packageJSON.i18n.outputLocales;
export const defaultLocale = packageJSON.i18n.entry;

export const NEXT_PUBLIC_DEFAULT_LANGUAGE =
  getEnv("NEXT_PUBLIC_DEFAULT_LANGUAGE") ??
  process.env.NEXT_PUBLIC_DEFAULT_LANGUAGE ??
  defaultLocale;

export const NEXT_PUBLIC_SITE_URL =
  getEnv("NEXT_PUBLIC_SITE_URL") ?? process.env.NEXT_PUBLIC_SITE_URL;
export const NEXT_PUBLIC_API_URL = getEnv("NEXT_PUBLIC_API_URL") ?? process.env.NEXT_PUBLIC_API_URL;
export const NEXT_PUBLIC_ADMIN_PATH =
  getEnv("NEXT_PUBLIC_ADMIN_PATH") ?? process.env.NEXT_PUBLIC_ADMIN_PATH ?? "/admin";

export const NEXT_PUBLIC_DEFAULT_USER_EMAIL =
  getEnv("NEXT_PUBLIC_DEFAULT_USER_EMAIL") ?? process.env.NEXT_PUBLIC_DEFAULT_USER_EMAIL;
export const NEXT_PUBLIC_DEFAULT_USER_PASSWORD =
  getEnv("NEXT_PUBLIC_DEFAULT_USER_PASSWORD") ?? process.env.NEXT_PUBLIC_DEFAULT_USER_PASSWORD;
