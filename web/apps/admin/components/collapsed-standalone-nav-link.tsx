"use client";

import { SidebarMenuButton } from "@workspace/ui/components/sidebar";
import type { ReactNode } from "react";
import { AdminLink } from "./admin-link";

type CollapsedStandaloneNavLinkProps = {
  href: string;
  label: string;
  isActive: boolean;
  icon: ReactNode;
};

export function CollapsedStandaloneNavLink({
  href,
  label,
  isActive,
  icon,
}: CollapsedStandaloneNavLinkProps) {
  return (
    <SidebarMenuButton
      asChild
      size="sm"
      className="h-8 justify-center"
      isActive={isActive}
      aria-label={label}
    >
      <AdminLink href={href}>{icon}</AdminLink>
    </SidebarMenuButton>
  );
}
