import { createTranslator } from "next-intl";
import { getLocale } from "@/utils/common";

export async function getTranslations(namespace: string) {
  const locale = getLocale();
  const messages = (await import(`./${locale}/${namespace}.json`)).default;
  return createTranslator({ locale, messages });
}
