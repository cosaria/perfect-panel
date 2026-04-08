import { forwardRef, type MouseEvent, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { stripAdminPath, toAdminPath } from "@/utils/admin-path";

interface LinkProps extends React.AnchorHTMLAttributes<HTMLAnchorElement> {
  href: string;
  children: ReactNode;
  replace?: boolean;
}

function isAbsoluteHref(href: string) {
  return /^https?:\/\//.test(href);
}

function isPlainLeftClick(event: MouseEvent<HTMLAnchorElement>) {
  return event.button === 0 && !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey;
}

export default forwardRef<HTMLAnchorElement, LinkProps>(function Link(
  { href, onClick, replace = false, target, download, ...props },
  ref,
) {
  const navigate = useNavigate();
  const resolvedHref = isAbsoluteHref(href) ? href : toAdminPath(href);

  return (
    <a
      {...props}
      ref={ref}
      href={resolvedHref}
      target={target}
      download={download}
      onClick={(event) => {
        onClick?.(event);

        if (
          event.defaultPrevented ||
          isAbsoluteHref(href) ||
          target === "_blank" ||
          download !== undefined ||
          !isPlainLeftClick(event)
        ) {
          return;
        }

        event.preventDefault();
        navigate(stripAdminPath(resolvedHref), { replace });
      }}
    />
  );
});
