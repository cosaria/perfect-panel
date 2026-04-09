import { locales, VITE_DEFAULT_LANGUAGE } from "@/config/constants";
import enMessages from "./en-US";
import zhMessages from "./zh-CN";

const allMessages: Record<string, Record<string, unknown>> = {
  "en-US": enMessages,
  "zh-CN": zhMessages,
};

function isSupportedLocale(locale?: string | null): locale is string {
  return Boolean(locale && locales.includes(locale));
}

export function getClientLocale(): string {
  if (typeof document !== "undefined") {
    const match = document.cookie.match(/(?:^|;\s*)locale=([^;]+)/);
    if (isSupportedLocale(match?.[1])) {
      return match[1];
    }
  }
  if (typeof window !== "undefined") {
    const storedLocale = window.localStorage?.getItem("locale");
    if (isSupportedLocale(storedLocale)) {
      return storedLocale;
    }
  }
  if (typeof navigator !== "undefined") {
    const browserLang = navigator.language?.split(",")?.[0] || "";
    if (isSupportedLocale(browserLang)) {
      return browserLang;
    }
  }
  return VITE_DEFAULT_LANGUAGE;
}

export function getMessages(locale: string): Record<string, unknown> {
  return allMessages[locale] || enMessages;
}
