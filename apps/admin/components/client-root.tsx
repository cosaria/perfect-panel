"use client";

import { Toaster } from "@workspace/ui/components/sonner";
import { getLangDir } from "@workspace/ui/hooks/use-lang-dir";
import { NextIntlClientProvider } from "next-intl";
import NextTopLoader from "nextjs-toploader";
import type React from "react";
import Providers from "@/components/providers";
import { getClientLocale, getMessages } from "@/locales/client";

export default function ClientRoot({ children }: { children: React.ReactNode }) {
  const locale = getClientLocale();
  const messages = getMessages(locale);

  return (
    <html suppressHydrationWarning lang={locale} dir={getLangDir(locale)}>
      <body
        suppressHydrationWarning
        className="size-full min-h-[calc(100dvh-env(safe-area-inset-top))] font-sans antialiased"
      >
        <NextIntlClientProvider locale={locale} messages={messages}>
          <NextTopLoader showSpinner={false} />
          <Providers>
            <Toaster richColors closeButton />
            {children}
          </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
