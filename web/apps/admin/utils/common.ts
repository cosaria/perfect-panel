import { isBrowser } from "@workspace/ui/utils";
import { intlFormat } from "date-fns";
import Cookies from "universal-cookie";
import { locales, VITE_ADMIN_PATH, VITE_DEFAULT_LANGUAGE } from "@/config/constants";
import { dispatchAdminLocaleChange } from "@/src/runtime-intl";
import { ADMIN_HOME_PATH, stripAdminPath, toAdminPath } from "./admin-path";

const cookies = new Cookies(null, {
  path: "/",
  maxAge: 365 * 24 * 60 * 60,
});

export function getLocale() {
  const browserLocale =
    typeof navigator !== "undefined" ? navigator.language?.split(",")?.[0] || "" : "";
  const defaultLocale = locales.includes(browserLocale) ? browserLocale : "";
  const cookies = new Cookies(null, { path: "/" });
  const cookieLocale = cookies.get("locale") || "";
  const storedLocale =
    typeof window !== "undefined" ? window.localStorage?.getItem("locale") || "" : "";
  const locale = cookieLocale || storedLocale || defaultLocale || VITE_DEFAULT_LANGUAGE;
  return locale;
}

export function setLocale(value: string) {
  cookies.set("locale", value);
  if (typeof window !== "undefined") {
    window.localStorage?.setItem("locale", value);
  }
  dispatchAdminLocaleChange(value);
}

export function setAuthorization(value: string) {
  return cookies.set("Authorization", value);
}

export function getAuthorization(value?: string) {
  const Authorization = isBrowser() ? cookies.get("Authorization") : value;
  if (!Authorization) return;
  return Authorization;
}

export function setRedirectUrl(value?: string) {
  if (value) {
    sessionStorage.setItem("redirect-url", value);
  }
}

export function getRedirectUrl() {
  const redirectUrl = sessionStorage.getItem("redirect-url") ?? ADMIN_HOME_PATH;
  return toAdminPath(stripAdminPath(redirectUrl));
}

export function Logout() {
  if (!isBrowser()) return;
  cookies.remove("Authorization");
  const pathname = location.pathname;
  if (!["", "/", VITE_ADMIN_PATH].includes(pathname)) {
    setRedirectUrl(location.pathname);
    location.href = toAdminPath("/");
  }
}

export function formatDate(date?: Date | number, showTime: boolean = true) {
  if (!date) return;
  const timeZone = localStorage.getItem("timezone") || "UTC";
  return intlFormat(date, {
    year: "numeric",
    month: "numeric",
    day: "numeric",
    ...(showTime && {
      hour: "numeric",
      minute: "numeric",
      second: "numeric",
    }),
    hour12: false,
    timeZone,
  });
}
