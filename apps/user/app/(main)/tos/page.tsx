import { Markdown } from "@workspace/ui/custom-components/markdown";
import { getTos } from "@/services/common-api/sdk.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";

export default async function Page() {
  const { data } = await getTos({
    baseUrl: (NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "") + "/v1/common",
  });
  return (
    <div className="container py-8">
      <Markdown>{data?.tos_content || ""}</Markdown>
    </div>
  );
}
