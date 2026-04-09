"use client";

import { Toaster } from "@workspace/ui/components/sonner";
import { getLangDir } from "@workspace/ui/hooks/use-lang-dir";
import type React from "react";
import { useEffect } from "react";
import Providers from "@/components/providers";
import { getClientLocale, getMessages } from "@/locales/client";
import { NextIntlClientProvider, useLocale } from "@/src/runtime-intl";
import RouterTopLoader from "@/src/router-top-loader";

function DocumentLocaleSync() {
  const locale = useLocale();

  useEffect(() => {
    document.documentElement.lang = locale;
    document.documentElement.dir = getLangDir(locale);
  }, [locale]);

  return null;
}

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
          <DocumentLocaleSync />
          <RouterTopLoader showSpinner={false} />
          <Providers>
            <Toaster richColors closeButton />
            {children}
          </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
