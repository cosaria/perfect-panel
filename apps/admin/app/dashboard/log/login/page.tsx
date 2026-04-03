"use client";

import { Badge } from "@workspace/ui/components/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@workspace/ui/components/tooltip";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { UserDetail } from "@/app/dashboard/user/user-detail";
import { IpLink } from "@/components/ip-link";
import { ProTable } from "@/components/pro-table";
import { filterLoginLog } from "@/services/admin-api/sdk.gen";
import type { LoginLog } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";

export default function LoginLogPage() {
  const t = useTranslations("log");
  const sp = useSearchParams();

  const today = new Date().toISOString().split("T")[0];

  const initialFilters = {
    date: sp.get("date") || today,
    user_id: sp.get("user_id") ? Number(sp.get("user_id")) : undefined,
  };
  return (
    <ProTable<LoginLog, { date?: string; user_id?: number }>
      header={{ title: t("title.login") }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: "user",
          header: t("column.user"),
          cell: ({ row }) => (
            <div>
              <Badge className="capitalize">{row.original.method}</Badge>{" "}
              <UserDetail id={Number(row.original.user_id)} />
            </div>
          ),
        },

        {
          accessorKey: "login_ip",
          header: t("column.ip"),
          cell: ({ row }) => <IpLink ip={String(row.original.login_ip || "")} />,
        },
        {
          accessorKey: "user_agent",
          header: t("column.userAgent"),
          cell: ({ row }) => {
            const userAgent = String(row.original.user_agent || "");
            return (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div className="max-w-48 cursor-help truncate">{userAgent}</div>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p className="max-w-md break-words">{userAgent}</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            );
          },
        },
        {
          accessorKey: "success",
          header: t("column.success"),
          cell: ({ row }) => (
            <Badge variant={row.original.success ? "default" : "destructive"}>
              {row.original.success ? t("success") : t("failed")}
            </Badge>
          ),
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
        const { data } = await filterLoginLog({
          query: {
            page: pagination.page,
            size: pagination.size,
            date: filter?.date || "",
            search: "",
            user_id: filter?.user_id ? Number(filter.user_id) : 0,
          },
        });
        const list = (data?.list || []) as LoginLog[];
        const total = Number(data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
