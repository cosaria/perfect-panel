import { createTranslator } from "@workspace/ui/components/i18n-provider";
import { getLocale } from "@/utils/common";

export async function getTranslations(namespace: string) {
  const locale = getLocale();
  const messages = (await import(`./${locale}/${namespace}.json`)).default;
  return createTranslator({ locale, messages });
}
