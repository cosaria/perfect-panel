import { Toaster } from "@workspace/ui/components/sonner";
import Providers from "@/components/providers";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import { client as adminClient } from "@/services/admin-api/client.gen";
import { currentUser } from "@/services/admin-api/sdk.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { getGlobalConfig } from "@/services/common-api/sdk.gen";
import type { GetGlobalConfigResponse } from "@/services/common-api/types.gen";
import type { User } from "@/services/admin-api/types.gen";
import "@workspace/ui/globals.css";
import { getLangDir } from "@workspace/ui/hooks/use-lang-dir";
import type { Metadata, Viewport } from "next";
import { unstable_noStore as noStore } from "next/cache";
// import { Geist, Geist_Mono } from 'next/font/google';
import { cookies } from "next/headers";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import { PublicEnvScript } from "next-runtime-env";
import NextTopLoader from "nextjs-toploader";
import type React from "react";

// const fontSans = Geist({
//   subsets: ['latin'],
//   variable: '--font-sans',
// });

// const fontMono = Geist_Mono({
//   subsets: ['latin'],
//   variable: '--font-mono',
// });

export async function generateMetadata(): Promise<Metadata> {
  noStore();

  const baseUrl = NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "";

  let site: GetGlobalConfigResponse["site"] | undefined;
  try {
    const { data } = await getGlobalConfig({
      client: commonClient,
      baseUrl: `${baseUrl}/v1/common`,
    });
    site = data?.site || undefined;
  } catch (error) {
    console.log("Error fetching global config:", error);
  }

  const defaultMetadata = {
    title: {
      default: site?.site_name || `PPanel`,
      template: `%s | ${site?.site_name || "PPanel"}`,
    },
    description: site?.site_desc || "",
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

  const baseUrl = NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "";

  let config: GetGlobalConfigResponse | undefined;
  let user: User | undefined;

  try {
    const { data } = await getGlobalConfig({
      client: commonClient,
      baseUrl: `${baseUrl}/v1/common`,
    });
    config = data;
  } catch (error) {
    console.log("Error fetching global config:", error);
  }

  if (Authorization) {
    try {
      const { data } = await currentUser({
        client: adminClient,
        baseUrl: `${baseUrl}/v1/admin`,
        headers: { Authorization },
      });
      if (data?.is_admin) {
        user = data;
      }
    } catch (error) {
      console.log("Error fetching current user:", error);
    }
  }

  return (
    <html suppressHydrationWarning lang={locale} dir={getLangDir(locale)}>
      <head>
        <PublicEnvScript />
      </head>
      <body
        suppressHydrationWarning
        //  ${fontSans.variable} ${fontMono.variable}
        className={`size-full min-h-[calc(100dvh-env(safe-area-inset-top))] font-sans antialiased`}
      >
        <NextIntlClientProvider messages={messages}>
          <NextTopLoader showSpinner={false} />
          <Providers common={{ ...config }} user={user}>
            <Toaster richColors closeButton />
            {children}
          </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
