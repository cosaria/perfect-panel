"use client";

import { SidebarInset, SidebarProvider } from "@workspace/ui/components/sidebar";
import Announcement from "@/components/announcement";
import { getAuthorization } from "@/utils/common";
import { SidebarLeft } from "./sidebar-left";
import { SidebarRight } from "./sidebar-right";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const auth = getAuthorization();

  return (
    <SidebarProvider className="container">
      <SidebarLeft className="sticky top-[84px] hidden w-52 border-r-0 bg-transparent lg:flex" />
      <SidebarInset className="relative p-4">{children}</SidebarInset>
      <SidebarRight className="sticky top-[84px] hidden w-52 border-r-0 bg-transparent 2xl:flex" />
      <Announcement type="popup" Authorization={auth} />
    </SidebarProvider>
  );
}
