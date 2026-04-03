import { Markdown } from "@workspace/ui/custom-components/markdown";
import { getPrivacyPolicy } from "@/services/common-api/sdk.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";

export default async function Page() {
  const { data } = await getPrivacyPolicy({
    baseUrl: (NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "") + "/v1/common",
  });
  return (
    <div className="container py-8">
      <Markdown>{data?.privacy_policy || ""}</Markdown>
    </div>
  );
}
