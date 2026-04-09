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
export const VITE_CDN_URL =
  getEnv("VITE_CDN_URL") || getBuildEnv("VITE_CDN_URL") || "https://cdn.jsdelivr.net";

export const VITE_DEFAULT_USER_EMAIL =
  getEnv("VITE_DEFAULT_USER_EMAIL") ?? getBuildEnv("VITE_DEFAULT_USER_EMAIL");
export const VITE_DEFAULT_USER_PASSWORD =
  getEnv("VITE_DEFAULT_USER_PASSWORD") ?? getBuildEnv("VITE_DEFAULT_USER_PASSWORD");

export const VITE_EMAIL = getEnv("VITE_EMAIL") ?? getBuildEnv("VITE_EMAIL");

export const VITE_TELEGRAM_LINK = getEnv("VITE_TELEGRAM_LINK") ?? getBuildEnv("VITE_TELEGRAM_LINK");
export const VITE_DISCORD_LINK = getEnv("VITE_DISCORD_LINK") ?? getBuildEnv("VITE_DISCORD_LINK");
export const VITE_GITHUB_LINK = getEnv("VITE_GITHUB_LINK") ?? getBuildEnv("VITE_GITHUB_LINK");
export const VITE_LINKEDIN_LINK = getEnv("VITE_LINKEDIN_LINK") ?? getBuildEnv("VITE_LINKEDIN_LINK");
export const VITE_TWITTER_LINK = getEnv("VITE_TWITTER_LINK") ?? getBuildEnv("VITE_TWITTER_LINK");
export const VITE_INSTAGRAM_LINK =
  getEnv("VITE_INSTAGRAM_LINK") ?? getBuildEnv("VITE_INSTAGRAM_LINK");

export const VITE_HOME_USER_COUNT = (() => {
  const value = getEnv("VITE_HOME_USER_COUNT") ?? getBuildEnv("VITE_HOME_USER_COUNT");
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return 999;
  return numberValue;
})();

export const VITE_HOME_SERVER_COUNT = (() => {
  const value = getEnv("VITE_HOME_SERVER_COUNT") ?? getBuildEnv("VITE_HOME_SERVER_COUNT");
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return 999;
  return numberValue;
})();

export const VITE_HOME_LOCATION_COUNT = (() => {
  const value = getEnv("VITE_HOME_LOCATION_COUNT") ?? getBuildEnv("VITE_HOME_LOCATION_COUNT");
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return 999;
  return numberValue;
})();
