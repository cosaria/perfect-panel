"use client";

import { type ComponentPropsWithoutRef, forwardRef, type MouseEvent } from "react";
import { useNavigate } from "react-router-dom";
import { stripAdminPath, toAdminPath } from "@/utils/admin-path";

export interface AdminLinkProps extends Omit<ComponentPropsWithoutRef<"a">, "href"> {
  href: string;
  replace?: boolean;
}

export const AdminLink = forwardRef<HTMLAnchorElement, AdminLinkProps>(function AdminLink(
  { href, replace = false, target, download, onClick, ...props },
  ref,
) {
  const navigate = useNavigate();
  const resolvedHref = toAdminPath(href);
  const navigationTarget = stripAdminPath(resolvedHref);
  const isInternalAdminHref = resolvedHref.startsWith("/") && target !== "_blank" && !download;

  function handleClick(event: MouseEvent<HTMLAnchorElement>) {
    onClick?.(event);

    if (
      event.defaultPrevented ||
      !isInternalAdminHref ||
      event.button !== 0 ||
      event.metaKey ||
      event.altKey ||
      event.ctrlKey ||
      event.shiftKey
    ) {
      return;
    }

    event.preventDefault();
    navigate(navigationTarget, { replace });
  }

  return (
    <a
      {...props}
      ref={ref}
      href={resolvedHref}
      target={target}
      download={download}
      onClick={handleClick}
    />
  );
});
