"use client";

import { SidebarInset, SidebarProvider } from "@workspace/ui/components/sidebar";
import { Header } from "@/components/header";
import { SidebarLeft } from "@/components/sidebar-left";

function getSidebarState(): boolean {
  if (typeof document === "undefined") return true;
  const match = document.cookie.match(/(?:^|;\s*)sidebar_state=([^;]+)/);
  return match?.[1] === "true";
}

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const defaultOpen = getSidebarState();
  return (
    <SidebarProvider defaultOpen={defaultOpen}>
      <SidebarLeft />
      <SidebarInset className="relative flex-grow overflow-hidden">
        <Header />
        <div className="h-[calc(100vh-56px)] flex-grow gap-4 overflow-auto p-4">{children}</div>
      </SidebarInset>
    </SidebarProvider>
  );
}
