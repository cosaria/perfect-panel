"use client";

import { Separator } from "@workspace/ui/components/separator";
import { Icon } from "@workspace/ui/custom-components/icon";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { Fragment } from "react";
import {
  NEXT_PUBLIC_DISCORD_LINK,
  NEXT_PUBLIC_EMAIL,
  NEXT_PUBLIC_GITHUB_LINK,
  NEXT_PUBLIC_INSTAGRAM_LINK,
  NEXT_PUBLIC_LINKEDIN_LINK,
  NEXT_PUBLIC_TELEGRAM_LINK,
  NEXT_PUBLIC_TWITTER_LINK,
} from "@/config/constants";
import useGlobalStore from "@/config/use-global";

const Links = [
  {
    icon: "uil:envelope",
    href: NEXT_PUBLIC_EMAIL ? `mailto:${NEXT_PUBLIC_EMAIL}` : undefined,
  },
  {
    icon: "uil:telegram",
    href: NEXT_PUBLIC_TELEGRAM_LINK,
  },
  {
    icon: "uil:twitter",
    href: NEXT_PUBLIC_TWITTER_LINK,
  },
  {
    icon: "uil:discord",
    href: NEXT_PUBLIC_DISCORD_LINK,
  },
  {
    icon: "uil:instagram",
    href: NEXT_PUBLIC_INSTAGRAM_LINK,
  },
  {
    icon: "uil:linkedin",
    href: NEXT_PUBLIC_LINKEDIN_LINK,
  },
  {
    icon: "uil:github",
    href: NEXT_PUBLIC_GITHUB_LINK,
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
