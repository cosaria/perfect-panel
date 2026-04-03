import { getLocale } from "next-intl/server";
import { getSubscription } from "@/services/user-api/sdk.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import { Content } from "./content";

export async function ProductShowcase() {
  try {
    const locale = await getLocale();
    const { data } = await getSubscription({
      query: { language: locale },
      baseUrl: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "",
    });
    const subscriptionList = data?.list || [];

    if (subscriptionList.length === 0) return null;

    return <Content subscriptionData={subscriptionList} />;
  } catch (_error) {
    return null;
  }
}
