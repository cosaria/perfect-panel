"use client";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@workspace/ui/components/sidebar";
import { Icon } from "@workspace/ui/custom-components/icon";
import Link from "@/src/compat/app-link";
import { usePathname } from "@/src/compat/app-navigation";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { navs } from "@/config/navs";

export function SidebarLeft({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const t = useTranslations("menu");
  const pathname = usePathname();
  return (
    <Sidebar collapsible="none" side="left" {...props}>
      <SidebarContent>
        <SidebarMenu>
          {navs.map((nav) => (
            <SidebarGroup key={nav.title}>
              {nav.items && <SidebarGroupLabel>{t(nav.title)}</SidebarGroupLabel>}
              <SidebarGroupContent>
                <SidebarMenu>
                  {(nav.items || [nav]).map((item) => (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton
                        asChild
                        tooltip={t(item.title)}
                        isActive={item.url === pathname}
                      >
                        <Link href={item.url}>
                          {item.icon && <Icon icon={item.icon} />}
                          <span>{t(item.title)}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          ))}
        </SidebarMenu>
      </SidebarContent>
    </Sidebar>
  );
}
