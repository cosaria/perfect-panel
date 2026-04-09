"use client";

import { createInstance, type i18n as I18nInstance } from "i18next";
import type React from "react";
import { createContext, useContext, useMemo } from "react";
import { I18nextProvider, initReactI18next, useTranslation } from "react-i18next";

type Messages = Record<string, unknown>;
type TranslatorOptions = Record<string, unknown>;
type Translator = (key: string, options?: TranslatorOptions) => string;
type ProviderProps = {
  children: React.ReactNode;
  locale: string;
  messages: Messages;
};
type LocaleContextValue = {
  locale: string;
  messages: Messages;
};

const LocaleContext = createContext<LocaleContextValue>({
  locale: "en-US",
  messages: {},
});

function createI18nInstance(locale: string, messages: Messages): I18nInstance {
  const instance = createInstance();

  instance.use(initReactI18next);
  void instance.init({
    fallbackLng: locale,
    initAsync: false,
    interpolation: {
      escapeValue: false,
    },
    lng: locale,
    resources: {
      [locale]: {
        translation: messages,
      },
    },
    returnNull: false,
  });

  return instance;
}

function joinTranslationKey(namespace: string | undefined, key: string) {
  if (!namespace) {
    return key;
  }

  if (!key) {
    return namespace;
  }

  return `${namespace}.${key}`;
}

export function createTranslator({
  locale,
  messages,
}: {
  locale: string;
  messages: Messages;
}): Translator {
  const i18n = createI18nInstance(locale, messages);

  return (key, options) => i18n.t(key, options);
}

export function NextIntlClientProvider({ children, locale, messages }: ProviderProps) {
  const i18n = useMemo(() => createI18nInstance(locale, messages), [locale, messages]);

  return (
    <LocaleContext.Provider value={{ locale, messages }}>
      <I18nextProvider i18n={i18n}>{children}</I18nextProvider>
    </LocaleContext.Provider>
  );
}

export function useLocale() {
  return useContext(LocaleContext).locale;
}

export function useTranslations(namespace?: string) {
  const { t } = useTranslation();

  return (key: string, options?: TranslatorOptions) =>
    t(joinTranslationKey(namespace, key), options);
}
