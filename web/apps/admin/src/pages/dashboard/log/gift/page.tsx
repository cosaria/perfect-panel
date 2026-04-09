"use client";

import { Badge } from "@workspace/ui/components/badge";
import { useSearchParams } from "@/utils/router";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { Suspense } from "react";
import { UserDetail, UserSubscribeDetail } from "@/src/pages/dashboard/user/user-detail";
import { Display } from "@/components/display";
import { OrderLink } from "@/components/order-link";
import { ProTable } from "@/components/pro-table";
import { filterGiftLog } from "@/services/admin-api/sdk.gen";
import type { GiftLog } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";

function GiftLogPageContent() {
  const t = useTranslations("log");
  const sp = useSearchParams();

  const today = new Date().toISOString().split("T")[0];

  const getGiftTypeText = (type: number) => {
    const typeText = t(`type.${type}`);
    if (typeText === `log.type.${type}`) {
      return `${t("unknown")} (${type})`;
    }
    return typeText;
  };

  const initialFilters = {
    date: sp.get("date") || today,
    user_id: sp.get("user_id") ? Number(sp.get("user_id")) : undefined,
  };
  return (
    <ProTable<GiftLog, { date?: string; user_id?: number }>
      header={{ title: t("title.gift") }}
      initialFilters={initialFilters}
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
          accessorKey: "order_no",
          header: t("column.orderNo"),
          cell: ({ row }) => <OrderLink orderId={row.original.order_no} />,
        },
        {
          accessorKey: "amount",
          header: t("column.amount"),
          cell: ({ row }) => <Display type="currency" value={row.original.amount} />,
        },
        {
          accessorKey: "balance",
          header: t("column.balance"),
          cell: ({ row }) => <Display type="currency" value={row.original.balance} />,
        },
        {
          accessorKey: "type",
          header: t("column.type"),
          cell: ({ row }) => <Badge>{getGiftTypeText(row.original.type)}</Badge>,
        },
        { accessorKey: "remark", header: t("column.remark") },
        {
          accessorKey: "timestamp",
          header: t("column.time"),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: "date", type: "date" },
        { key: "user_id", placeholder: t("column.userId") },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterGiftLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            user_id: filter?.user_id ? Number(filter.user_id) : 0,
          },
        });
        const list = (data?.list || []) as GiftLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}

export default function GiftLogPage() {
  return (
    <Suspense>
      <GiftLogPageContent />
    </Suspense>
  );
}
