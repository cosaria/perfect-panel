import { Markdown } from "@workspace/ui/custom-components/markdown";
import { getTos } from "@/services/common/common";

export default async function Page() {
  const { data } = await getTos();
  return (
    <div className="container py-8">
      <Markdown>{data.data?.tos_content || ""}</Markdown>
    </div>
  );
}
