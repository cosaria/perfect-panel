"use client";

import { useEffect, useState } from "react";
import { getSubscription } from "@/services/user-api/sdk.gen";
import type { Subscribe } from "@/services/user-api/types.gen";
import { getClientLocale } from "@/locales/client";
import { Content } from "./content";

export function ProductShowcase() {
  const [subscriptionList, setSubscriptionList] = useState<Subscribe[]>([]);

  useEffect(() => {
    const locale = getClientLocale();
    getSubscription({
      query: { language: locale },
    })
      .then(({ data }) => {
        setSubscriptionList(data?.list || []);
      })
      .catch(() => {});
  }, []);

  if (subscriptionList.length === 0) return null;

  return <Content subscriptionData={subscriptionList} />;
}
