"use client";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@workspace/ui/components/select";
import { Icon } from "@workspace/ui/custom-components/icon";
import { getCountry } from "@workspace/ui/utils";
import { useLocale } from "@workspace/ui/components/i18n-provider";
import { locales } from "@/config/constants";
import { setLocale } from "@/utils/common";

const languages = {
  "en-US": "English",
  "zh-CN": "简体中文",
} as const;

export default function LanguageSwitch() {
  const locale = useLocale();
  const country = getCountry(locale);

  const handleLanguageChange = (value: string) => {
    setLocale(value);
  };

  return (
    <Select value={locale} onValueChange={handleLanguageChange}>
      <SelectTrigger className="hover:bg-accent hover:text-accent-foreground w-auto border-none bg-transparent p-2 shadow-none focus:ring-0 [&>svg]:hidden">
        <SelectValue>
          <div className="flex items-center">
            <Icon icon={`flagpack:${country?.alpha2.toLowerCase()}`} className="!size-5" />
            <span className="sr-only">{languages[locale as keyof typeof languages]}</span>
          </div>
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {locales.map(getCountry).map((item) => (
          <SelectItem key={`${item?.lang}-${item?.alpha2}`} value={`${item?.lang}-${item?.alpha2}`}>
            <div className="flex items-center gap-2">
              <Icon icon={`flagpack:${item?.alpha2.toLowerCase()}`} className="!size-5" />
              {languages[`${item?.lang}-${item?.alpha2}` as keyof typeof languages]}
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
