"use client";

import "@/utils/setup-clients";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryStreamedHydration } from "@tanstack/react-query-next-experimental";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import type React from "react";
import { useEffect, useState } from "react";
import useGlobalStore from "@/config/use-global";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import { getGlobalConfig } from "@/services/common-api/sdk.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { queryUserInfo } from "@/services/user-api/sdk.gen";
import { getAuthorization, Logout } from "@/utils/common";
import Loading from "./loading";

export default function Providers({
  children,
}: {
  children: React.ReactNode;
}) {
  const [loading, setLoading] = useState(true);
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
  const [customHtml, setCustomHtml] = useState("");

  useEffect(() => {
    const baseUrl = NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "";

    const init = async () => {
      try {
        const { data: config } = await getGlobalConfig({
          client: commonClient,
          baseUrl: `${baseUrl}/v1/common`,
        });
        if (config) {
          setCommon(config);
          if (config.site?.site_name) {
            document.title = config.site.site_name;
          }
          if (config.site?.custom_html) {
            setCustomHtml(config.site.custom_html);
          }
        }

        const auth = getAuthorization();
        if (auth) {
          try {
            const { data: user } = await queryUserInfo({
              baseUrl: baseUrl,
              headers: { Authorization: auth },
            });
            if (user) {
              setUser(user);
            } else {
              Logout();
            }
          } catch {
            Logout();
          }
        }
      } finally {
        setTimeout(() => {
          setLoading(false);
        }, 1000);
      }
    };

    init();
  }, [setCommon, setUser]);

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const invite = searchParams.get("invite");
    if (invite) {
      localStorage.setItem("invite", invite);
    }
  }, []);

  return (
    <NextThemesProvider attribute="class" defaultTheme="system" enableSystem>
      <QueryClientProvider client={queryClient}>
        <ReactQueryStreamedHydration>
          <Loading loading={loading || queryClient.isMutating() > 0} />
          {children}
        </ReactQueryStreamedHydration>
      </QueryClientProvider>
      {customHtml && (
        <div id="custom_html" dangerouslySetInnerHTML={{ __html: customHtml }} />
      )}
    </NextThemesProvider>
  );
}
