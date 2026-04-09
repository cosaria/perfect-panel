"use client";

import { Badge } from "@workspace/ui/components/badge";
import { Button } from "@workspace/ui/components/button";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@workspace/ui/components/sheet";
import { Switch } from "@workspace/ui/components/switch";
import { ConfirmButton } from "@workspace/ui/custom-components/confirm-button";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { type ReactNode, useState } from "react";
import { toast } from "sonner";
import { IpLink } from "@/components/ip-link";
import { ProTable } from "@/components/pro-table";
import { getUserSubscribeDevices, kickOfflineByUserDevice } from "@/services/admin-api/sdk.gen";
import type { UserDevice } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";

export function SubscriptionDetail({
  trigger,
  userId,
  subscriptionId,
}: {
  trigger: ReactNode;
  userId: number;
  subscriptionId: number;
}) {
  const t = useTranslations("user");
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>{trigger}</SheetTrigger>
      <SheetContent side="right" className="w-[700px] max-w-full md:max-w-screen-md">
        <SheetHeader>
          <SheetTitle>{t("onlineDevices")}</SheetTitle>
        </SheetHeader>
        <div className="mt-4 max-h-[calc(100dvh-120px)] overflow-y-auto">
          <ProTable<UserDevice, Record<string, unknown>>
            columns={[
              {
                accessorKey: "enabled",
                header: t("enable"),
                cell: ({ row }) => (
                  <Switch
                    checked={row.getValue("enabled")}
                    onChange={(checked) => {
                      console.log("Switch:", checked);
                    }}
                  />
                ),
              },
              { accessorKey: "id", header: "ID" },
              { accessorKey: "identifier", header: "IMEI" },
              { accessorKey: "user_agent", header: t("userAgent") },
              {
                accessorKey: "ip",
                header: "IP",
                cell: ({ row }) => <IpLink ip={row.getValue("ip")} />,
              },
              {
                accessorKey: "online",
                header: t("loginStatus"),
                cell: ({ row }) => (
                  <Badge variant={row.getValue("online") ? "default" : "destructive"}>
                    {row.getValue("online") ? t("online") : t("offline")}
                  </Badge>
                ),
              },
              {
                accessorKey: "updated_at",
                header: t("lastSeen"),
                cell: ({ row }) => formatDate(row.getValue("updated_at")),
              },
            ]}
            request={async (pagination) => {
              const { data } = await getUserSubscribeDevices({
                query: {
                  user_id: userId,
                  subscribe_id: subscriptionId,
                  ...pagination,
                },
              });
              return {
                list: data?.list || [],
                total: data?.total || 0,
              };
            }}
            actions={{
              render: (row) => {
                if (!row.identifier) return [];
                return [
                  <ConfirmButton
                    key="offline"
                    trigger={<Button variant="destructive">{t("confirmOffline")}</Button>}
                    title={t("confirmOffline")}
                    description={t("kickOfflineConfirm", { ip: row.ip })}
                    onConfirm={async () => {
                      await kickOfflineByUserDevice({ body: { id: row.id } });
                      toast.success(t("kickOfflineSuccess"));
                    }}
                    cancelText={t("cancel")}
                    confirmText={t("confirm")}
                  />,
                ];
              },
            }}
          />
        </div>
      </SheetContent>
    </Sheet>
  );
}
