"use client";

import { Badge } from "@workspace/ui/components/badge";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { UserDetail } from "@/app/dashboard/user/user-detail";
import { Display } from "@/components/display";
import { OrderLink } from "@/components/order-link";
import { ProTable } from "@/components/pro-table";
import { filterCommissionLog } from "@/services/admin-api/sdk.gen";
import type { CommissionLog } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";

export default function CommissionLogPage() {
  const t = useTranslations("log");
  const sp = useSearchParams();

  const today = new Date().toISOString().split("T")[0];

  const getCommissionTypeText = (type: number) => {
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
    <ProTable<CommissionLog, { date?: string; user_id?: number }>
      header={{ title: t("title.commission") }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: "user",
          header: t("column.user"),
          cell: ({ row }) => <UserDetail id={Number(row.original.user_id)} />,
        },
        {
          accessorKey: "amount",
          header: t("column.amount"),
          cell: ({ row }) => <Display type="currency" value={row.original.amount} />,
        },
        {
          accessorKey: "order_no",
          header: t("column.orderNo"),
          cell: ({ row }) => <OrderLink orderId={row.original.order_no} />,
        },
        {
          accessorKey: "type",
          header: t("column.type"),
          cell: ({ row }) => <Badge>{getCommissionTypeText(row.original.type)}</Badge>,
        },
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
        const { data } = await filterCommissionLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            user_id: filter?.user_id ? Number(filter.user_id) : 0,
          },
        });
        const list = (data?.list || []) as CommissionLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
