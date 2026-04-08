"use client";

import { Badge } from "@workspace/ui/components/badge";
import { Button } from "@workspace/ui/components/button";
import { formatBytes } from "@workspace/ui/utils";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { Suspense } from "react";
import { ProTable } from "@/components/pro-table";
import { filterServerTrafficLog } from "@/services/admin-api/sdk.gen";
import type { ServerTrafficLog } from "@/services/admin-api/types.gen";
import { useServer } from "@/store/server";

function ServerTrafficLogPageContent() {
  const t = useTranslations("log");
  const sp = useSearchParams();
  const { getServerName } = useServer();

  const today = new Date().toISOString().split("T")[0];

  const initialFilters = {
    date: sp.get("date") || today,
    server_id: sp.get("server_id") ? Number(sp.get("server_id")) : undefined,
  };
  return (
    <ProTable<ServerTrafficLog, { date?: string; server_id?: number }>
      header={{ title: t("title.serverTraffic") }}
      initialFilters={initialFilters}
      actions={{
        render: (row) => [
          <Button key="detail" asChild>
            <Link
              href={`/dashboard/log/traffic-details?date=${row.date}&server_id=${row.server_id}`}
            >
              {t("detail")}
            </Link>
          </Button>,
        ],
      }}
      columns={[
        {
          accessorKey: "server_id",
          header: t("column.server"),
          cell: ({ row }) => {
            return (
              <div className="flex items-center gap-2">
                <Badge>{row.original.server_id}</Badge>
                <span>{getServerName(row.original.server_id)}</span>
              </div>
            );
          },
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
        { accessorKey: "date", header: t("column.date") },
      ]}
      params={[
        { key: "date", type: "date" },
        { key: "server_id", placeholder: t("column.serverId") },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterServerTrafficLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            server_id: filter?.server_id ? Number(filter.server_id) : 0,
          },
        });
        const list = (data?.list || []) as ServerTrafficLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}

export default function ServerTrafficLogPage() {
  return (
    <Suspense>
      <ServerTrafficLogPageContent />
    </Suspense>
  );
}
