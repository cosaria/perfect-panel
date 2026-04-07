"use client";

import { Card } from "@workspace/ui/components/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@workspace/ui/components/dialog";
import { Icon } from "@workspace/ui/custom-components/icon";
import { Markdown } from "@workspace/ui/custom-components/markdown";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";
import { queryAnnouncement } from "@/services/user-api/sdk.gen";
import type { Announcement as AnnouncementType } from "@/services/user-api/types.gen";
import { Empty } from "../empty";

export default function Announcement({
  type,
  Authorization,
}: {
  type: "popup" | "pinned";
  Authorization?: string;
}) {
  const t = useTranslations("dashboard");
  const [data, setData] = useState<AnnouncementType | undefined>();

  useEffect(() => {
    queryAnnouncement({
      body: {
        page: 1,
        size: 10,
        pinned: type === "pinned",
        popup: type === "popup",
      },
      headers: Authorization ? { Authorization } : undefined,
    })
      .then(({ data: result }) => {
        setData(result?.announcements?.find((item) => item[type]) ?? undefined);
      })
      .catch(() => {});
  }, [type, Authorization]);

  if (!data) return null;

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
