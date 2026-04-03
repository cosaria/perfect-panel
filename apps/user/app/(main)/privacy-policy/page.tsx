import { Markdown } from "@workspace/ui/custom-components/markdown";
import { getPrivacyPolicy } from "@/services/common/common";

export default async function Page() {
  const { data } = await getPrivacyPolicy();
  return (
    <div className="container py-8">
      <Markdown>{data.data?.privacy_policy || ""}</Markdown>
    </div>
  );
}
