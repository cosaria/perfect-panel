"use client";

import { formatBytes } from "@workspace/ui/utils";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { UserDetail, UserSubscribeDetail } from "@/app/dashboard/user/user-detail";
import { ProTable } from "@/components/pro-table";
import { filterTrafficLogDetails } from "@/services/admin-api/sdk.gen";
import type { TrafficLogDetails } from "@/services/admin-api/types.gen";
import { useServer } from "@/store/server";
import { formatDate } from "@/utils/common";

export default function TrafficDetailsPage() {
  const t = useTranslations("log");
  const sp = useSearchParams();
  const { getServerName } = useServer();

  const today = new Date().toISOString().split("T")[0];

  const initialFilters = {
    date: sp.get("date") || today,
    server_id: sp.get("server_id") ? Number(sp.get("server_id")) : undefined,
    user_id: sp.get("user_id") ? Number(sp.get("user_id")) : undefined,
    subscribe_id: sp.get("subscribe_id") ? Number(sp.get("subscribe_id")) : undefined,
  };
  return (
    <ProTable<
      TrafficLogDetails,
      {
        date?: string;
        server_id?: number;
        user_id?: number;
        subscribe_id?: number;
      }
    >
      header={{ title: t("title.trafficDetails") }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: "server_id",
          header: t("column.server"),
          cell: ({ row }) => (
            <span>
              {getServerName(row.original.server_id)} ({row.original.server_id})
            </span>
          ),
        },
        {
          accessorKey: "user_id",
          header: t("column.user"),
          cell: ({ row }) => <UserDetail id={Number(row.original.user_id)} />,
        },
        {
          accessorKey: "subscribe_id",
          header: t("column.subscribe"),
          cell: ({ row }) => (
            <UserSubscribeDetail id={Number(row.original.subscribe_id)} enabled hoverCard />
          ),
        },
        {
          accessorKey: "upload",
          header: t("column.upload"),
          cell: ({ row }) => formatBytes(row.original.upload),
        },
        {
          accessorKey: "download",
          header: t("column.download"),
          cell: ({ row }) => formatBytes(row.original.download),
        },
        {
          accessorKey: "timestamp",
          header: t("column.time"),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: "date", type: "date" },
        { key: "server_id", placeholder: t("column.serverId") },
        { key: "user_id", placeholder: t("column.userId") },
        { key: "subscribe_id", placeholder: t("column.subscribeId") },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterTrafficLogDetails({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            server_id: filter?.server_id ? Number(filter.server_id) : 0,
            user_id: filter?.user_id ? Number(filter.user_id) : 0,
            subscribe_id: filter?.subscribe_id ? Number(filter.subscribe_id) : 0,
          },
        });
        const list = (data?.list || []) as TrafficLogDetails[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
