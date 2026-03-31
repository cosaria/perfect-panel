import Providers from '@/components/providers';
import { currentUser } from '@/services/admin/user';
import { getGlobalConfig } from '@/services/common/common';
import { Toaster } from '@workspace/ui/components/sonner';
import '@workspace/ui/globals.css';
import { getLangDir } from '@workspace/ui/hooks/use-lang-dir';
import { Metadata, Viewport } from 'next';
import { NextIntlClientProvider } from 'next-intl';
import { getLocale, getMessages } from 'next-intl/server';
import { PublicEnvScript } from 'next-runtime-env';
import { unstable_noStore as noStore } from 'next/cache';
// import { Geist, Geist_Mono } from 'next/font/google';
import { cookies } from 'next/headers';
import NextTopLoader from 'nextjs-toploader';
import React from 'react';

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

  let site: API.SiteConfig | undefined;
  await getGlobalConfig({ skipErrorHandler: true })
    .then((res) => {
      const config = res.data.data;
      site = config?.site || undefined;
    })
    .catch((error) => {
      console.log('Error fetching global config:', error);
    });

  const defaultMetadata = {
    title: {
      default: site?.site_name || `PPanel`,
      template: `%s | ${site?.site_name || 'PPanel'}`,
    },
    description: site?.site_desc || '',
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
      user = await currentUser({
        skipErrorHandler: true,
        Authorization,
      }).then((res) => {
        if (res.data.data?.is_admin) {
          return res.data.data;
        }
        return undefined;
      });
    } catch (error) {
      console.log('Error fetching current user:', error);
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
