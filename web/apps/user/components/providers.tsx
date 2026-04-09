"use client";

import "@/utils/setup-clients";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@workspace/ui/components/theme-provider";
import type React from "react";
import { useEffect, useState } from "react";
import useGlobalStore from "@/config/use-global";
import { client as commonClient } from "@/services/common-api/client.gen";
import { getGlobalConfig } from "@/services/common-api/sdk.gen";
import { queryUserInfo } from "@/services/user-api/sdk.gen";
import { getAuthorization, Logout } from "@/utils/common";
import Loading from "./loading";
import { TrustedCustomHtml } from "./trusted-custom-html";

export default function Providers({ children }: { children: React.ReactNode }) {
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
    const init = async () => {
      try {
        const { data: config } = await getGlobalConfig({
          client: commonClient,
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
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <QueryClientProvider client={queryClient}>
        <Loading loading={loading || queryClient.isMutating() > 0} />
        {children}
      </QueryClientProvider>
      {customHtml && <TrustedCustomHtml html={customHtml} />}
    </ThemeProvider>
  );
}
