import { Toaster } from "@workspace/ui/components/sonner";
import { RawHtml } from "@workspace/ui/custom-components/raw-html";
import Providers from "@/components/providers";
import { getGlobalConfig } from "@/services/common-api/sdk.gen";
import { queryUserInfo } from "@/services/user-api/sdk.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import type { GetGlobalConfigResponse } from "@/services/common-api/types.gen";
import type { User } from "@/services/user-api/types.gen";
import "@workspace/ui/globals.css";
import { getLangDir } from "@workspace/ui/hooks/use-lang-dir";
import { unstable_noStore as noStore } from "next/cache";
// import { Geist, Geist_Mono } from "next/font/google";
import { cookies } from "next/headers";
import type { Metadata, Viewport } from "next/types";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import { PublicEnvScript } from "next-runtime-env";
import NextTopLoader from "nextjs-toploader";
import type React from "react";

// const geistSans = Geist({
//   variable: "--font-geist-sans",
//   subsets: ["latin"],
// });

// const geistMono = Geist_Mono({
//   variable: "--font-geist-mono",
//   subsets: ["latin"],
// });

export async function generateMetadata(): Promise<Metadata> {
  noStore();
  let site: GetGlobalConfigResponse["site"] | undefined;

  const ssrBaseUrl = (NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "") + "/v1/common";
  try {
    const { data: config } = await getGlobalConfig({
      baseUrl: ssrBaseUrl,
    });
    site = config?.site || undefined;
  } catch (error) {
    console.error("Error fetching global config:", error);
  }

  const defaultMetadata = {
    title: {
      default: site?.site_name || `PPanel`,
      template: `%s | ${site?.site_name || "PPanel"}`,
    },
    description: site?.site_desc || "",
    keywords: site?.keywords || "",
    icons: {
      icon: site?.site_logo
        ? [
            {
              url: site.site_logo,
              sizes: "any",
            },
          ]
        : [
            { url: "/favicon.ico", sizes: "48x48" },
            { url: "/favicon.svg", type: "image/svg+xml" },
          ],
      apple: site?.site_logo || "/apple-touch-icon.png",
    },
    manifest: "/site.webmanifest",
  };
  return defaultMetadata;
}

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#FFFFFF" },
    { media: "(prefers-color-scheme: dark)", color: "#000000" },
  ],
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const Authorization = (await cookies()).get("Authorization")?.value;

  const locale = await getLocale();
  const messages = await getMessages();

  let config: GetGlobalConfigResponse | undefined;
  let user: User | undefined;

  const ssrBaseUrl = (NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "") + "/v1/common";
  try {
    const { data } = await getGlobalConfig({
      baseUrl: ssrBaseUrl,
    });
    config = data;
  } catch (error) {
    console.log("Error fetching global config:", error);
  }

  if (Authorization) {
    try {
      const { data } = await queryUserInfo({
        baseUrl: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "",
        headers: {
          Authorization,
        },
      });
      user = data;
    } catch (error) {
      console.log("Error fetching user info:", error);
    }
  }

  const customHtml = config?.site.custom_html || "";

  return (
    <html suppressHydrationWarning lang={locale} dir={getLangDir(locale)}>
      <head>
        <PublicEnvScript />
      </head>
      <body
        suppressHydrationWarning
        // ${geistSans.variable} ${geistMono.variable}
        className={`size-full min-h-[calc(100dvh-env(safe-area-inset-top))] font-sans antialiased`}
      >
        <NextIntlClientProvider messages={messages}>
          <NextTopLoader showSpinner={false} />
          <Providers common={{ ...config }} user={user}>
            <Toaster richColors closeButton />
            {children}
          </Providers>
        </NextIntlClientProvider>
        <RawHtml id="custom_html" html={customHtml} />
      </body>
    </html>
  );
}
