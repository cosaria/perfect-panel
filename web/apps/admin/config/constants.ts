import packageJSON from "../package.json";

declare global {
  interface Window {
    __ENV?: Record<string, string>;
  }
}

type PublicImportMetaEnv = {
  [key: string]: string | undefined;
};

function getEnv(key: string): string | undefined {
  if (typeof window !== "undefined") {
    return window.__ENV?.[key];
  }
  return undefined;
}

function getBuildEnv(key: string): string | undefined {
  return (import.meta as ImportMeta & { env?: PublicImportMetaEnv }).env?.[key] ?? process.env[key];
}

export const locales = packageJSON.i18n.outputLocales;
export const defaultLocale = packageJSON.i18n.entry;

export const VITE_DEFAULT_LANGUAGE =
  getEnv("VITE_DEFAULT_LANGUAGE") ?? getBuildEnv("VITE_DEFAULT_LANGUAGE") ?? defaultLocale;

export const VITE_SITE_URL = getEnv("VITE_SITE_URL") ?? getBuildEnv("VITE_SITE_URL");
export const VITE_API_URL = getEnv("VITE_API_URL") ?? getBuildEnv("VITE_API_URL");
export const VITE_ADMIN_PATH =
  getEnv("VITE_ADMIN_PATH") ?? getBuildEnv("VITE_ADMIN_PATH") ?? "/admin";

export const VITE_DEFAULT_USER_EMAIL =
  getEnv("VITE_DEFAULT_USER_EMAIL") ?? getBuildEnv("VITE_DEFAULT_USER_EMAIL");
export const VITE_DEFAULT_USER_PASSWORD =
  getEnv("VITE_DEFAULT_USER_PASSWORD") ?? getBuildEnv("VITE_DEFAULT_USER_PASSWORD");

// Transitional aliases keep existing imports working while Unit 1 migrates callers.
export const NEXT_PUBLIC_DEFAULT_LANGUAGE = VITE_DEFAULT_LANGUAGE;
export const NEXT_PUBLIC_SITE_URL = VITE_SITE_URL;
export const NEXT_PUBLIC_API_URL = VITE_API_URL;
export const NEXT_PUBLIC_ADMIN_PATH = VITE_ADMIN_PATH;
export const NEXT_PUBLIC_DEFAULT_USER_EMAIL = VITE_DEFAULT_USER_EMAIL;
export const NEXT_PUBLIC_DEFAULT_USER_PASSWORD = VITE_DEFAULT_USER_PASSWORD;
