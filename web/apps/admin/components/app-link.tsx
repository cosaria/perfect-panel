"use client";

import { forwardRef, type MouseEvent, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { stripAdminPath, toAdminPath } from "@/utils/admin-path";

interface LinkProps extends React.AnchorHTMLAttributes<HTMLAnchorElement> {
  children: ReactNode;
  href: string;
  replace?: boolean;
}

function isAbsoluteHref(href: string) {
  return /^https?:\/\//.test(href);
}

export default forwardRef<HTMLAnchorElement, LinkProps>(function AppLink(
  { href, replace = false, target, download, onClick, ...props },
  ref,
) {
  const navigate = useNavigate();
  const resolvedHref = isAbsoluteHref(href) ? href : toAdminPath(href);
  const navigationTarget = isAbsoluteHref(href) ? href : stripAdminPath(resolvedHref);
  const isInternalAdminHref =
    !isAbsoluteHref(href) && target !== "_blank" && download === undefined;

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
