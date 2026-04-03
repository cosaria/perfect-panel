"use client";

import { Badge } from "@workspace/ui/components/badge";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { UserDetail, UserSubscribeDetail } from "@/app/dashboard/user/user-detail";
import { OrderLink } from "@/components/order-link";
import { ProTable } from "@/components/pro-table";
import { filterResetSubscribeLog } from "@/services/admin-api/sdk.gen";
import type { ResetSubscribeLog } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";

export default function ResetSubscribeLogPage() {
  const t = useTranslations("log");
  const sp = useSearchParams();

  const today = new Date().toISOString().split("T")[0];

  const getResetSubscribeTypeText = (type: number) => {
    const typeText = t(`type.${type}`);
    if (typeText === `log.type.${type}`) {
      return `${t("unknown")} (${type})`;
    }
    return typeText;
  };

  const initialFilters = {
    date: sp.get("date") || today,
    user_subscribe_id: sp.get("user_subscribe_id")
      ? Number(sp.get("user_subscribe_id"))
      : undefined,
  };
  return (
    <ProTable<ResetSubscribeLog, { date?: string; user_subscribe_id?: number }>
      header={{ title: t("title.resetSubscribe") }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: "user",
          header: t("column.user"),
          cell: ({ row }) => <UserDetail id={Number(row.original.user_id)} />,
        },
        {
          accessorKey: "user_subscribe_id",
          header: t("column.subscribeId"),
          cell: ({ row }) => (
            <UserSubscribeDetail id={Number(row.original.user_subscribe_id)} enabled hoverCard />
          ),
        },
        {
          accessorKey: "type",
          header: t("column.type"),
          cell: ({ row }) => <Badge>{getResetSubscribeTypeText(row.original.type)}</Badge>,
        },
        {
          accessorKey: "order_no",
          header: t("column.orderNo"),
          cell: ({ row }) => <OrderLink orderId={row.original.order_no} />,
        },
        {
          accessorKey: "timestamp",
          header: t("column.time"),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: "date", type: "date" },
        { key: "user_subscribe_id", placeholder: t("column.subscribeId") },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterResetSubscribeLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            user_subscribe_id: filter?.user_subscribe_id ? Number(filter.user_subscribe_id) : 0,
          },
        });
        const list = (data?.list || []) as ResetSubscribeLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
