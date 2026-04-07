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
export const NEXT_PUBLIC_CDN_URL =
  getEnv("NEXT_PUBLIC_CDN_URL") || process.env.NEXT_PUBLIC_CDN_URL || "https://cdn.jsdelivr.net";

export const NEXT_PUBLIC_DEFAULT_USER_EMAIL =
  getEnv("NEXT_PUBLIC_DEFAULT_USER_EMAIL") ?? process.env.NEXT_PUBLIC_DEFAULT_USER_EMAIL;
export const NEXT_PUBLIC_DEFAULT_USER_PASSWORD =
  getEnv("NEXT_PUBLIC_DEFAULT_USER_PASSWORD") ?? process.env.NEXT_PUBLIC_DEFAULT_USER_PASSWORD;

export const NEXT_PUBLIC_EMAIL = getEnv("NEXT_PUBLIC_EMAIL") ?? process.env.NEXT_PUBLIC_EMAIL;

export const NEXT_PUBLIC_TELEGRAM_LINK =
  getEnv("NEXT_PUBLIC_TELEGRAM_LINK") ?? process.env.NEXT_PUBLIC_TELEGRAM_LINK;
export const NEXT_PUBLIC_DISCORD_LINK =
  getEnv("NEXT_PUBLIC_DISCORD_LINK") ?? process.env.NEXT_PUBLIC_DISCORD_LINK;
export const NEXT_PUBLIC_GITHUB_LINK =
  getEnv("NEXT_PUBLIC_GITHUB_LINK") ?? process.env.NEXT_PUBLIC_GITHUB_LINK;
export const NEXT_PUBLIC_LINKEDIN_LINK =
  getEnv("NEXT_PUBLIC_LINKEDIN_LINK") ?? process.env.NEXT_PUBLIC_LINKEDIN_LINK;
export const NEXT_PUBLIC_TWITTER_LINK =
  getEnv("NEXT_PUBLIC_TWITTER_LINK") ?? process.env.NEXT_PUBLIC_TWITTER_LINK;
export const NEXT_PUBLIC_INSTAGRAM_LINK =
  getEnv("NEXT_PUBLIC_INSTAGRAM_LINK") ?? process.env.NEXT_PUBLIC_INSTAGRAM_LINK;

export const NEXT_PUBLIC_HOME_USER_COUNT = (() => {
  const value = getEnv("NEXT_PUBLIC_HOME_USER_COUNT") ?? process.env.NEXT_PUBLIC_HOME_USER_COUNT;
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return 999;
  return numberValue;
})();

export const NEXT_PUBLIC_HOME_SERVER_COUNT = (() => {
  const value =
    getEnv("NEXT_PUBLIC_HOME_SERVER_COUNT") ?? process.env.NEXT_PUBLIC_HOME_SERVER_COUNT;
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return 999;
  return numberValue;
})();

export const NEXT_PUBLIC_HOME_LOCATION_COUNT = (() => {
  const value =
    getEnv("NEXT_PUBLIC_HOME_LOCATION_COUNT") ?? process.env.NEXT_PUBLIC_HOME_LOCATION_COUNT;
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return 999;
  return numberValue;
})();

export const NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT =
  getEnv("NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT") ??
  process.env.NEXT_PUBLIC_HIDDEN_TUTORIAL_DOCUMENT;
