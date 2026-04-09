"use client";

import type React from "react";
import { createContext, useContext, useEffect, useState } from "react";
import {
  applyResolvedTheme,
  DEFAULT_THEME,
  DEFAULT_THEME_STORAGE_KEY,
  readStoredTheme,
  resolveTheme,
  type ResolvedTheme,
  type Theme,
} from "../lib/theme.js";

type ThemeProviderProps = {
  attribute?: "class";
  children: React.ReactNode;
  defaultTheme?: Theme;
  enableSystem?: boolean;
  storageKey?: string;
};

type ThemeContextValue = {
  resolvedTheme: ResolvedTheme;
  setTheme: (theme: Theme) => void;
  systemTheme: ResolvedTheme;
  theme: Theme;
};

type MediaQueryLike = {
  addEventListener?: (type: "change", listener: (event: { matches: boolean }) => void) => void;
  addListener?: (listener: (event: { matches: boolean }) => void) => void;
  matches: boolean;
  removeEventListener?: (type: "change", listener: (event: { matches: boolean }) => void) => void;
  removeListener?: (listener: (event: { matches: boolean }) => void) => void;
};

const DEFAULT_RESOLVED_THEME: ResolvedTheme = "light";

const ThemeContext = createContext<ThemeContextValue>({
  resolvedTheme: DEFAULT_RESOLVED_THEME,
  setTheme() {},
  systemTheme: DEFAULT_RESOLVED_THEME,
  theme: DEFAULT_THEME,
});

function getSystemTheme(): ResolvedTheme {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return DEFAULT_RESOLVED_THEME;
  }

  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}

export function ThemeProvider({
  attribute = "class",
  children,
  defaultTheme = DEFAULT_THEME,
  enableSystem = true,
  storageKey = DEFAULT_THEME_STORAGE_KEY,
}: ThemeProviderProps) {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window === "undefined") {
      return defaultTheme;
    }

    return readStoredTheme(window.localStorage, storageKey, defaultTheme);
  });
  const [systemTheme, setSystemTheme] = useState<ResolvedTheme>(() => getSystemTheme());

  useEffect(() => {
    if (!enableSystem || typeof window === "undefined" || typeof window.matchMedia !== "function") {
      return;
    }

    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)") as MediaQueryLike;
    const handleChange = (event: { matches: boolean }) => {
      setSystemTheme(event.matches ? "dark" : "light");
    };

    handleChange(mediaQuery);

    if (typeof mediaQuery.addEventListener === "function") {
      mediaQuery.addEventListener("change", handleChange);
      return () => mediaQuery.removeEventListener?.("change", handleChange);
    }

    mediaQuery.addListener?.(handleChange);
    return () => mediaQuery.removeListener?.(handleChange);
  }, [enableSystem]);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    window.localStorage?.setItem(storageKey, theme);
  }, [storageKey, theme]);

  const resolvedTheme =
    enableSystem && theme === "system"
      ? resolveTheme(theme, systemTheme)
      : theme === "dark"
        ? "dark"
        : "light";

  useEffect(() => {
    if (attribute !== "class" || typeof document === "undefined") {
      return;
    }

    applyResolvedTheme(document.documentElement, resolvedTheme);
  }, [attribute, resolvedTheme]);

  return (
    <ThemeContext.Provider
      value={{
        resolvedTheme,
        setTheme: setThemeState,
        systemTheme,
        theme,
      }}
    >
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  return useContext(ThemeContext);
}
