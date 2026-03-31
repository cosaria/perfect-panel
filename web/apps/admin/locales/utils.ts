import { getLocale } from '@/utils/common';
import { createTranslator } from 'next-intl';

export async function getTranslations(namespace: string) {
  const locale = getLocale();
  const messages = (await import(`./${locale}/${namespace}.json`)).default;
  return createTranslator({ locale, messages });
}
