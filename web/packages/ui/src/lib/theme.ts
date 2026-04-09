export type Theme = "dark" | "light" | "system";
export type ResolvedTheme = Exclude<Theme, "system">;

type ClassListLike = {
  add: (...tokens: string[]) => void;
  remove: (...tokens: string[]) => void;
};

type StorageLike = {
  getItem: (key: string) => null | string;
};

type ThemeTarget = {
  classList: ClassListLike;
  style?: {
    colorScheme?: string;
  };
};

export const DEFAULT_THEME: Theme = "system";
export const DEFAULT_THEME_STORAGE_KEY = "theme";
export const THEME_CLASS_NAMES: ResolvedTheme[] = ["light", "dark"];

export function isTheme(value: null | string | undefined): value is Theme {
  return value === "light" || value === "dark" || value === "system";
}

export function readStoredTheme(
  storage: null | StorageLike | undefined,
  key: string,
  fallbackTheme: Theme = DEFAULT_THEME,
): Theme {
  const storedTheme = storage?.getItem(key);
  return isTheme(storedTheme) ? storedTheme : fallbackTheme;
}

export function resolveTheme(theme: Theme, systemTheme: ResolvedTheme): ResolvedTheme {
  return theme === "system" ? systemTheme : theme;
}

export function resolveSonnerTheme(theme: Theme, resolvedTheme: ResolvedTheme): ResolvedTheme {
  return theme === "system" ? resolvedTheme : theme;
}

export function applyResolvedTheme(target: ThemeTarget, resolvedTheme: ResolvedTheme) {
  target.classList.remove(...THEME_CLASS_NAMES);
  target.classList.add(resolvedTheme);
  if (target.style) {
    target.style.colorScheme = resolvedTheme;
  }
}
