import {
  NextIntlClientProvider as BaseNextIntlClientProvider,
  createTranslator,
  useLocale,
  useTranslations,
} from "next-intl";
import { type ComponentProps, useEffect, useState } from "react";
import { getClientLocale, getMessages } from "@/locales/client";

export const USER_LOCALE_CHANGE_EVENT = "ppanel:user-locale-change";

type IntlProviderProps = ComponentProps<typeof BaseNextIntlClientProvider>;
type LocaleChangeDetail = {
  locale?: string;
};

function resolveRuntimeLocale(detail?: LocaleChangeDetail) {
  const locale = detail?.locale || getClientLocale();

  return {
    locale,
    messages: getMessages(locale),
  };
}

export function dispatchUserLocaleChange(locale: string) {
  if (typeof window === "undefined") {
    return;
  }

  window.dispatchEvent(
    new CustomEvent<LocaleChangeDetail>(USER_LOCALE_CHANGE_EVENT, {
      detail: { locale },
    }),
  );
}

export function NextIntlClientProvider({
  children,
  locale,
  messages,
  ...props
}: IntlProviderProps) {
  const [runtimeIntlState, setRuntimeIntlState] = useState(() => ({
    locale,
    messages,
  }));

  useEffect(() => {
    setRuntimeIntlState({ locale, messages });
  }, [locale, messages]);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    const handleLocaleChange = (event: Event) => {
      setRuntimeIntlState(resolveRuntimeLocale((event as CustomEvent<LocaleChangeDetail>).detail));
    };

    window.addEventListener(USER_LOCALE_CHANGE_EVENT, handleLocaleChange);

    return () => {
      window.removeEventListener(USER_LOCALE_CHANGE_EVENT, handleLocaleChange);
    };
  }, []);

  return (
    <BaseNextIntlClientProvider
      {...props}
      locale={runtimeIntlState.locale}
      messages={runtimeIntlState.messages}
    >
      {children}
    </BaseNextIntlClientProvider>
  );
}

export { createTranslator, useLocale, useTranslations };
