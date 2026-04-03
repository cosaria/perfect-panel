import { getLocale } from "next-intl/server";
import { getSubscription } from "@/services/user-api/sdk.gen";
import type { Subscribe } from "@/services/user-api/types.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import Content from "./content";

export default async function Page({
  searchParams,
}: {
  searchParams: Promise<{
    id: string;
  }>;
}) {
  const { id } = await searchParams;
  const locale = await getLocale();
  let subscriptionList: Subscribe[] = [];
  try {
    const { data } = await getSubscription({
      query: { language: locale },
      baseUrl: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "",
    });
    subscriptionList = data?.list || [];
  } catch {
    // silently handle SSR errors
  }
  const subscription = subscriptionList.find((item) => item.id === Number(id));

  return (
    <main className="container space-y-16">
      <Content subscription={subscription} />
    </main>
  );
}
