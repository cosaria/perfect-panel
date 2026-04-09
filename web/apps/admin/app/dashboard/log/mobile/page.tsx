"use client";

import { Badge } from "@workspace/ui/components/badge";
import { useSearchParams } from "@/utils/router";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { Suspense } from "react";
import { ProTable } from "@/components/pro-table";
import { filterMobileLog } from "@/services/admin-api/sdk.gen";
import type { MessageLog } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";

function MobileLogPageContent() {
  const t = useTranslations("log");
  const sp = useSearchParams();

  const today = new Date().toISOString().split("T")[0];

  const initialFilters = {
    search: sp.get("search") || undefined,
    date: sp.get("date") || today,
  };
  return (
    <ProTable<MessageLog, { search?: string; date?: string }>
      header={{ title: t("title.mobile") }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: "platform",
          header: t("column.platform"),
          cell: ({ row }) => <Badge>{row.getValue("platform")}</Badge>,
        },
        { accessorKey: "to", header: t("column.to") },
        { accessorKey: "subject", header: t("column.subject") },
        {
          accessorKey: "content",
          header: t("column.content"),
          cell: ({ row }) => (
            <pre className="max-w-[480px] overflow-auto whitespace-pre-wrap break-words text-xs">
              {JSON.stringify(row.original.content || {}, null, 2)}
            </pre>
          ),
        },
        {
          accessorKey: "status",
          header: t("column.status"),
          cell: ({ row }) => {
            const status = row.original.status;
            const getStatusVariant = (status: number | undefined) => {
              if (status === 1) {
                return "default";
              } else if (status === 0) {
                return "destructive";
              }
              return "outline";
            };

            const getStatusText = (status: number | undefined) => {
              if (status === 1) return t("sent");
              if (status === 0) return t("failed");
              return t("unknown");
            };

            return <Badge variant={getStatusVariant(status)}>{getStatusText(status)}</Badge>;
          },
        },
        {
          accessorKey: "created_at",
          header: t("column.time"),
          cell: ({ row }) => formatDate(row.original.created_at),
        },
      ]}
      params={[{ key: "search" }, { key: "date", type: "date" }]}
      request={async (pagination, filter) => {
        const { data } = await filterMobileLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            search: filter?.search || "",
            date: filter?.date || "",
          },
        });
        const list = (data?.list || []) as MessageLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}

export default function MobileLogPage() {
  return (
    <Suspense>
      <MobileLogPageContent />
    </Suspense>
  );
}
