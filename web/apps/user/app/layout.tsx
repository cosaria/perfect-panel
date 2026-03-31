import Providers from '@/components/providers';
import { getGlobalConfig } from '@/services/common/common';
import { queryUserInfo } from '@/services/user/user';
import { Toaster } from '@workspace/ui/components/sonner';
import '@workspace/ui/globals.css';
import { getLangDir } from '@workspace/ui/hooks/use-lang-dir';
import { NextIntlClientProvider } from 'next-intl';
import { getLocale, getMessages } from 'next-intl/server';
import { PublicEnvScript } from 'next-runtime-env';
import { unstable_noStore as noStore } from 'next/cache';
// import { Geist, Geist_Mono } from "next/font/google";
import { cookies } from 'next/headers';
import { Metadata, Viewport } from 'next/types';
import NextTopLoader from 'nextjs-toploader';
import React from 'react';

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
  let site: API.SiteConfig | undefined;

  await getGlobalConfig({ skipErrorHandler: true })
    .then((res) => {
      const config = res.data.data;
      site = config?.site || undefined;
    })
    .catch((error) => {
      console.error('Error fetching global config:', error);
    });

  const defaultMetadata = {
    title: {
      default: site?.site_name || `PPanel`,
      template: `%s | ${site?.site_name || 'PPanel'}`,
    },
    description: site?.site_desc || '',
    keywords: site?.keywords || '',
    icons: {
      icon: site?.site_logo
        ? [
            {
              url: site.site_logo,
              sizes: 'any',
            },
          ]
        : [
            { url: '/favicon.ico', sizes: '48x48' },
            { url: '/favicon.svg', type: 'image/svg+xml' },
          ],
      apple: site?.site_logo || '/apple-touch-icon.png',
    },
    manifest: '/site.webmanifest',
  };
  return defaultMetadata;
}

export const viewport: Viewport = {
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#FFFFFF' },
    { media: '(prefers-color-scheme: dark)', color: '#000000' },
  ],
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const Authorization = (await cookies()).get('Authorization')?.value;

  const locale = await getLocale();
  const messages = await getMessages();

  let config, user;

  try {
    config = await getGlobalConfig({ skipErrorHandler: true }).then((res) => res.data.data);
  } catch (error) {
    console.log('Error fetching global config:', error);
  }

  if (Authorization) {
    try {
      user = await queryUserInfo({
        skipErrorHandler: true,
        Authorization,
      }).then((res) => res.data.data);
    } catch (error) {
      console.log('Error fetching user info:', error);
    }
  }

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
        <div
          id='custom_html'
          dangerouslySetInnerHTML={{ __html: config?.site.custom_html || '' }}
        />
      </body>
    </html>
  );
}
