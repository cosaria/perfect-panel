"use client";

import { Button } from "@workspace/ui/components/button";
import { formatBytes } from "@workspace/ui/utils";
import { useSearchParams } from "@/utils/router";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { Suspense } from "react";
import { UserDetail, UserSubscribeDetail } from "@/app/dashboard/user/user-detail";
import { AdminLink } from "@/components/admin-link";
import { ProTable } from "@/components/pro-table";
import { filterUserSubscribeTrafficLog } from "@/services/admin-api/sdk.gen";
import type { UserSubscribeTrafficLog } from "@/services/admin-api/types.gen";

function SubscribeTrafficLogPageContent() {
  const t = useTranslations("log");
  const sp = useSearchParams();

  const today = new Date().toISOString().split("T")[0];

  const initialFilters = {
    date: sp.get("date") || today,
    user_id: sp.get("user_id") ? Number(sp.get("user_id")) : undefined,
    user_subscribe_id: sp.get("user_subscribe_id")
      ? Number(sp.get("user_subscribe_id"))
      : undefined,
  };
  return (
    <ProTable<
      UserSubscribeTrafficLog,
      { date?: string; user_id?: number; user_subscribe_id?: number }
    >
      header={{ title: t("title.subscribeTraffic") }}
      initialFilters={initialFilters}
      actions={{
        render: (row) => [
          <Button key="detail" asChild>
            <AdminLink
              href={`/dashboard/log/traffic-details?date=${row.date}&user_id=${row.user_id}&subscribe_id=${row.subscribe_id}`}
            >
              {t("detail")}
            </AdminLink>
          </Button>,
        ],
      }}
      columns={[
        {
          accessorKey: "user",
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
          accessorKey: "total",
          header: t("column.total"),
          cell: ({ row }) => formatBytes(row.original.total),
        },
        {
          accessorKey: "date",
          header: t("column.date"),
        },
      ]}
      params={[
        { key: "date", type: "date" },
        { key: "user_id", placeholder: t("column.userId") },
        { key: "user_subscribe_id", placeholder: t("column.subscribeId") },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterUserSubscribeTrafficLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            user_id: filter?.user_id ? Number(filter.user_id) : 0,
            user_subscribe_id: filter?.user_subscribe_id ? Number(filter.user_subscribe_id) : 0,
          },
        });
        const list = (data?.list || []) as UserSubscribeTrafficLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}

export default function SubscribeTrafficLogPage() {
  return (
    <Suspense>
      <SubscribeTrafficLogPageContent />
    </Suspense>
  );
}
