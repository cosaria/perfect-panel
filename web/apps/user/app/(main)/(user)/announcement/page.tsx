"use client";

import { useQuery } from "@tanstack/react-query";
import { Timeline } from "@workspace/ui/components/timeline";
import { Markdown } from "@workspace/ui/custom-components/markdown";
import { Empty } from "@/components/empty";
import { queryAnnouncement } from "@/services/user-api/sdk.gen";

export default function Page() {
  const { data } = useQuery({
    queryKey: ["queryAnnouncement"],
    queryFn: async () => {
      const { data } = await queryAnnouncement({
        body: {
          page: 1,
          size: 99,
          pinned: false,
          popup: false,
        },
      });
      return data?.announcements || [];
    },
  });
  return data && data.length > 0 ? (
    <Timeline
      data={
        data.map((item) => ({
          title: item.title,
          content: <Markdown>{item.content}</Markdown>,
        })) || []
      }
    />
  ) : (
    <Empty border />
  );
}
