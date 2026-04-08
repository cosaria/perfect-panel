import type { ReactNode } from "react";
import { useMemo } from "react";
import {
  useLocation,
  useNavigate,
  useParams as useReactRouterParams,
  useSearchParams as useReactRouterSearchParams,
} from "react-router-dom";

function isAbsoluteHref(href: string) {
  return /^https?:\/\//.test(href);
}

export function useRouter() {
  const navigate = useNavigate();

  return useMemo(
    () => ({
      push(href: string) {
        if (isAbsoluteHref(href)) {
          window.location.assign(href);
          return;
        }
        navigate(href);
      },
      replace(href: string) {
        if (isAbsoluteHref(href)) {
          window.location.replace(href);
          return;
        }
        navigate(href, { replace: true });
      },
      refresh() {
        window.location.reload();
      },
      back() {
        window.history.back();
      },
      forward() {
        window.history.forward();
      },
      prefetch() {
        return Promise.resolve();
      },
    }),
    [navigate],
  );
}

export function usePathname() {
  useLocation();
  return window.location.pathname;
}

export function useSearchParams() {
  const [searchParams] = useReactRouterSearchParams();
  return searchParams;
}

export function useParams<T extends Record<string, string | undefined>>() {
  return useReactRouterParams() as T;
}

export function useServerInsertedHTML(_callback: () => ReactNode) {
  // Vite user app 只运行在浏览器，不需要 Next server insert 生命周期。
}
