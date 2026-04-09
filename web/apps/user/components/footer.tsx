"use client";

import { Separator } from "@workspace/ui/components/separator";
import { Icon } from "@workspace/ui/custom-components/icon";
import Link from "@/components/app-link";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { Fragment } from "react";
import {
  VITE_DISCORD_LINK,
  VITE_EMAIL,
  VITE_GITHUB_LINK,
  VITE_INSTAGRAM_LINK,
  VITE_LINKEDIN_LINK,
  VITE_TELEGRAM_LINK,
  VITE_TWITTER_LINK,
} from "@/config/constants";
import useGlobalStore from "@/config/use-global";

const Links = [
  {
    icon: "uil:envelope",
    href: VITE_EMAIL ? `mailto:${VITE_EMAIL}` : undefined,
  },
  {
    icon: "uil:telegram",
    href: VITE_TELEGRAM_LINK,
  },
  {
    icon: "uil:twitter",
    href: VITE_TWITTER_LINK,
  },
  {
    icon: "uil:discord",
    href: VITE_DISCORD_LINK,
  },
  {
    icon: "uil:instagram",
    href: VITE_INSTAGRAM_LINK,
  },
  {
    icon: "uil:linkedin",
    href: VITE_LINKEDIN_LINK,
  },
  {
    icon: "uil:github",
    href: VITE_GITHUB_LINK,
  },
];

export default function Footer() {
  const { common } = useGlobalStore();
  const { site } = common;
  const t = useTranslations("auth");
  const activeLinks = Links.filter(
    (item): item is typeof item & { href: string } =>
      typeof item.href === "string" && item.href.length > 0,
  );

  return (
    <footer>
      <Separator className="my-14" />
      <div className="text-muted-foreground container mb-14 flex flex-wrap justify-between gap-4 text-sm">
        <nav className="flex flex-wrap items-center gap-2">
          {activeLinks.map((item, index) => (
            <Fragment key={item.href}>
              {index !== 0 && <Separator orientation="vertical" />}
              <Link href={item.href}>
                <Icon icon={item.icon} className="text-foreground size-5" />
              </Link>
            </Fragment>
          ))}
        </nav>
        <div>
          <strong className="text-foreground">{site.site_name}</strong> © All rights reserved.
          <div>
            <Link href="/tos" className="underline">
              {t("tos")}
            </Link>
            <Link href="/privacy-policy" className="ml-2 underline">
              {t("privacyPolicy")}
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
