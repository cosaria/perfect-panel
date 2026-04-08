"use client";

import { useRouter } from "next/navigation";
import { type ComponentPropsWithoutRef, forwardRef, type MouseEvent } from "react";
import { stripAdminPath, toAdminPath } from "@/utils/admin-path";

export interface AdminLinkProps extends Omit<ComponentPropsWithoutRef<"a">, "href"> {
  href: string;
  replace?: boolean;
}

function isPlainLeftClick(event: MouseEvent<HTMLAnchorElement>) {
  return event.button === 0 && !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey;
}

export const AdminLink = forwardRef<HTMLAnchorElement, AdminLinkProps>(function AdminLink(
  { href, onClick, replace = false, target, download, ...props },
  ref,
) {
  const router = useRouter();
  const resolvedHref = toAdminPath(href);
  const navigationTarget = stripAdminPath(resolvedHref);
  const isInternalAdminHref = resolvedHref.startsWith("/");

  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    onClick?.(event);

    if (
      event.defaultPrevented ||
      !isInternalAdminHref ||
      target === "_blank" ||
      download !== undefined ||
      !isPlainLeftClick(event)
    ) {
      return;
    }

    event.preventDefault();

    if (replace) {
      router.replace(navigationTarget);
      return;
    }

    router.push(navigationTarget);
  };

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
