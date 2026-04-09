"use client";

import "@/utils/setup-clients";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@workspace/ui/components/theme-provider";
import { usePathname } from "@/utils/router";
import type React from "react";
import { useEffect, useState } from "react";
import useGlobalStore from "@/config/use-global";
import { client as adminClient } from "@/services/admin-api/client.gen";
import { currentUser } from "@/services/admin-api/sdk.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { getGlobalConfig } from "@/services/common-api/sdk.gen";
import { useStatsStore } from "@/store/stats";
import { canonicalizeAdminBrowserPath } from "@/utils/admin-path";
import { getAuthorization, Logout } from "@/utils/common";

export default function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 5 * 1000,
            retry: false,
          },
        },
      }),
  );

  const { setCommon, setUser } = useGlobalStore();
  const pathname = usePathname();

  useEffect(() => {
    getGlobalConfig({
      client: commonClient,
    })
      .then(({ data }) => {
        if (data) {
          setCommon(data);
          if (data.site?.site_name) {
            document.title = data.site.site_name;
          }
        }
      })
      .catch(console.error);

    const auth = getAuthorization();
    if (auth) {
      currentUser({
        client: adminClient,
        headers: { Authorization: auth },
      })
        .then(({ data }) => {
          if (data?.is_admin) {
            setUser(data);
          } else {
            Logout();
          }
        })
        .catch(() => Logout());
    }
  }, [setCommon, setUser]);

  const { stats } = useStatsStore();

  useEffect(() => {
    stats();
  }, [stats]);

  useEffect(() => {
    if (typeof pathname !== "string") {
      return;
    }

    const currentPathname = window.location.pathname;
    const nextPathname = canonicalizeAdminBrowserPath(currentPathname);
    if (nextPathname === currentPathname) {
      return;
    }

    const nextUrl = `${nextPathname}${window.location.search}${window.location.hash}`;
    window.history.replaceState(window.history.state, "", nextUrl);
  }, [pathname]);

  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </ThemeProvider>
  );
}
