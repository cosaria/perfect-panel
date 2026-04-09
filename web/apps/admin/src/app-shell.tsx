"use client";

import { Toaster } from "@workspace/ui/components/sonner";
import { getLangDir } from "@workspace/ui/hooks/use-lang-dir";
import { type ReactNode, useEffect } from "react";
import Providers from "@/components/providers";
import { getClientLocale, getMessages } from "@/locales/client";
import { NextIntlClientProvider, useLocale } from "./compat/runtime-intl";
import RouterTopLoader from "./compat/router-top-loader";

function DocumentLocaleSync() {
  const locale = useLocale();

  useEffect(() => {
    document.documentElement.lang = locale;
    document.documentElement.dir = getLangDir(locale);
  }, [locale]);

  return null;
}

export default function AppShell({ children }: { children: ReactNode }) {
  const locale = getClientLocale();
  const messages = getMessages(locale);

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <DocumentLocaleSync />
      <RouterTopLoader showSpinner={false} />
      <Providers>
        <Toaster richColors closeButton />
        {children}
      </Providers>
    </NextIntlClientProvider>
  );
}
