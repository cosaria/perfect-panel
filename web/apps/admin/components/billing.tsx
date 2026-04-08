"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@workspace/ui/components/avatar";
import { Card, CardDescription, CardHeader, CardTitle } from "@workspace/ui/components/card";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

interface BillingProps {
  type: "dashboard" | "payment";
}

interface ItemType {
  logo: string;
  title: string;
  description: string;
  expiryDate: string;
  href: string;
}

async function fetchBillingList(type: string): Promise<ItemType[]> {
  try {
    let url: string;
    try {
      const response = await fetch(
        "https://api.github.com/repos/perfect-panel/ppanel-assets/commits",
      );
      const json = await response.json();
      const version = json[0]?.sha || "latest";
      url = `https://cdn.jsdmirror.com/gh/perfect-panel/ppanel-assets@${version}/billing/index.json`;
    } catch {
      url = "https://cdn.jsdmirror.com/gh/perfect-panel/ppanel-assets/billing/index.json";
    }

    const response = await fetch(url, {
      headers: { Accept: "application/json" },
    });
    const data = await response.json();
    const now = Date.now();

    return Array.isArray(data[type])
      ? data[type].filter((item: { expiryDate: string }) => {
          const expiryDate = Date.parse(item.expiryDate);
          return !Number.isNaN(expiryDate) && expiryDate > now;
        })
      : [];
  } catch {
    return [];
  }
}

export default function Billing({ type }: BillingProps) {
  const t = useTranslations("common.billing");
  const [list, setList] = useState<ItemType[]>([]);

  useEffect(() => {
    fetchBillingList(type).then(setList);
  }, [type]);

  if (!list.length) return null;

  return (
    <>
      <h1 className="text mt-2 font-bold">
        <span>{t("title")}</span>
        <span className="text-muted-foreground ml-2 text-xs">{t("description")}</span>
      </h1>
      <div className="grid gap-3 md:grid-cols-3 lg:grid-cols-6">
        {list.map((item) => (
          <Link href={item.href} target="_blank" key={item.href}>
            <Card className="h-full cursor-pointer">
              <CardHeader className="flex flex-row gap-2 p-3">
                <Avatar>
                  <AvatarImage src={item.logo} />
                  <AvatarFallback>{item.title}</AvatarFallback>
                </Avatar>
                <div>
                  <CardTitle>{item.title}</CardTitle>
                  <CardDescription className="mt-2">{item.description}</CardDescription>
                </div>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </div>
    </>
  );
}
