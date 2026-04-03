import { Card } from "@workspace/ui/components/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@workspace/ui/components/dialog";
import { Icon } from "@workspace/ui/custom-components/icon";
import { Markdown } from "@workspace/ui/custom-components/markdown";
import { getTranslations } from "next-intl/server";
import { queryAnnouncement } from "@/services/user-api/sdk.gen";
import type { Announcement as AnnouncementType } from "@/services/user-api/types.gen";
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from "@/config/constants";
import { Empty } from "../empty";

export default async function Announcement({
  type,
  Authorization,
}: {
  type: "popup" | "pinned";
  Authorization?: string;
}) {
  let data: AnnouncementType | undefined;
  try {
    const { data: result } = await queryAnnouncement({
      body: {
        page: 1,
        size: 10,
        pinned: type === "pinned",
        popup: type === "popup",
      },
      baseUrl: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL || "",
      headers: Authorization ? { Authorization } : undefined,
    });
    data = result?.announcements?.find((item) => item[type]) ?? undefined;
  } catch (_error) {
    /* empty */
  }
  if (!data) return null;

  const t = await getTranslations("dashboard");

  if (type === "popup") {
    return (
      <Dialog defaultOpen={!!data}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{data?.title}</DialogTitle>
          </DialogHeader>
          <Markdown>{data?.content}</Markdown>
        </DialogContent>
      </Dialog>
    );
  }
  if (type === "pinned") {
    return (
      <>
        <h2 className="flex items-center gap-1.5 font-semibold">
          <Icon icon="uil:bell" className="size-5" />
          {t("latestAnnouncement")}
        </h2>
        <Card className="p-6">
          {data?.content ? <Markdown>{data?.content}</Markdown> : <Empty />}
        </Card>
      </>
    );
  }
}
