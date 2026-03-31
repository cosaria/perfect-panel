import { getLocale } from "next-intl/server";
import { getSubscription } from "@/services/user/portal";
import { Content } from "./content";

export async function ProductShowcase() {
  try {
    const locale = await getLocale();
    const { data } = await getSubscription(
      {
        language: locale,
      },
      {
        skipErrorHandler: true,
      },
    );
    const subscriptionList = data.data?.list || [];

    if (subscriptionList.length === 0) return null;

    return <Content subscriptionData={subscriptionList} />;
  } catch (_error) {
    return null;
  }
}
