import { useMemo } from "react";
import {
  useLocation,
  useNavigate,
  useParams as useReactRouterParams,
  useSearchParams as useReactRouterSearchParams,
} from "react-router-dom";
import { stripAdminPath, toAdminPath } from "@/utils/admin-path";

function isAbsoluteHref(href: string) {
  return /^https?:\/\//.test(href);
}

function toInternalPath(href: string) {
  return stripAdminPath(toAdminPath(href));
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
        navigate(toInternalPath(href));
      },
      replace(href: string) {
        if (isAbsoluteHref(href)) {
          window.location.replace(href);
          return;
        }
        navigate(toInternalPath(href), { replace: true });
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
