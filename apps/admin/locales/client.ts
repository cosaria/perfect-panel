import { locales, NEXT_PUBLIC_DEFAULT_LANGUAGE } from "@/config/constants";
import enMessages from "./en-US";
import zhMessages from "./zh-CN";

const allMessages: Record<string, Record<string, unknown>> = {
  "en-US": enMessages,
  "zh-CN": zhMessages,
};

export function getClientLocale(): string {
  if (typeof document !== "undefined") {
    const match = document.cookie.match(/(?:^|;\s*)locale=([^;]+)/);
    if (match?.[1] && locales.includes(match[1])) {
      return match[1];
    }
  }
  if (typeof navigator !== "undefined") {
    const browserLang = navigator.language?.split(",")?.[0] || "";
    if (locales.includes(browserLang)) {
      return browserLang;
    }
  }
  return NEXT_PUBLIC_DEFAULT_LANGUAGE;
}

export function getMessages(locale: string): Record<string, unknown> {
  return allMessages[locale] || enMessages;
}
