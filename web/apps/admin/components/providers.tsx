"use client";

import "@/utils/setup-clients";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryStreamedHydration } from "@tanstack/react-query-next-experimental";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import type React from "react";
import { useEffect, useState } from "react";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import useGlobalStore from "@/config/use-global";
import { client as adminClient } from "@/services/admin-api/client.gen";
import { currentUser } from "@/services/admin-api/sdk.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { getGlobalConfig } from "@/services/common-api/sdk.gen";
import { useStatsStore } from "@/store/stats";
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

  useEffect(() => {
    const baseUrl = NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "";

    getGlobalConfig({
      client: commonClient,
      baseUrl: `${baseUrl}/v1/common`,
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
        baseUrl: `${baseUrl}/v1/admin`,
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

  return (
    <NextThemesProvider attribute="class" defaultTheme="system" enableSystem>
      <QueryClientProvider client={queryClient}>
        <ReactQueryStreamedHydration>{children}</ReactQueryStreamedHydration>
      </QueryClientProvider>
    </NextThemesProvider>
  );
}
