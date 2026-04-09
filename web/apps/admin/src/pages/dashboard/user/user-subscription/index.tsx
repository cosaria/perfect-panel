"use client";

import { Button } from "@workspace/ui/components/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@workspace/ui/components/dropdown-menu";
import { ConfirmButton } from "@workspace/ui/custom-components/confirm-button";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { useRef, useState } from "react";
import { toast } from "sonner";
import { AdminLink } from "@/components/admin-link";
import { Display } from "@/components/display";
import { ProTable, type ProTableActions } from "@/components/pro-table";
import useGlobalStore from "@/config/use-global";
import {
  createUserSubscribe,
  deleteUserSubscribe,
  getUserSubscribe,
  updateUserSubscribe,
} from "@/services/admin-api/sdk.gen";
import type { UserSubscribe } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";
import { SubscriptionDetail } from "./subscription-detail";
import { SubscriptionForm } from "./subscription-form";

export default function UserSubscription({ userId }: { userId: number }) {
  const t = useTranslations("user");
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);
  const { getUserSubscribeUrls } = useGlobalStore();

  return (
    <ProTable<UserSubscribe, Record<string, unknown>>
      action={ref}
      header={{
        title: t("subscriptionList"),
        toolbar: (
          <SubscriptionForm
            key="create"
            trigger={t("add")}
            title={t("createSubscription")}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              await createUserSubscribe({
                body: {
                  user_id: Number(userId),
                  subscribe_id: values.subscribe_id ?? 0,
                  traffic: values.traffic ?? 0,
                  expired_at: values.expired_at ?? 0,
                },
              });
              toast.success(t("createSuccess"));
              ref.current?.refresh();
              setLoading(false);
              return true;
            }}
          />
        ),
      }}
      columns={[
        {
          accessorKey: "id",
          header: "ID",
        },
        {
          accessorKey: "name",
          header: t("subscriptionName"),
          cell: ({ row }) => row.original.subscribe.name,
        },
        {
          accessorKey: "upload",
          header: t("upload"),
          cell: ({ row }) => <Display type="traffic" value={row.getValue("upload")} />,
        },
        {
          accessorKey: "download",
          header: t("download"),
          cell: ({ row }) => <Display type="traffic" value={row.getValue("download")} />,
        },
        {
          accessorKey: "traffic",
          header: t("totalTraffic"),
          cell: ({ row }) => <Display type="traffic" value={row.getValue("traffic")} unlimited />,
        },
        {
          accessorKey: "speed_limit",
          header: t("speedLimit"),
          cell: ({ row }) => {
            const speed = row.original?.subscribe?.speed_limit;
            return <Display type="trafficSpeed" value={speed} />;
          },
        },
        {
          accessorKey: "device_limit",
          header: t("deviceLimit"),
          cell: ({ row }) => {
            const limit = row.original?.subscribe?.device_limit;
            return <Display type="number" value={limit} unlimited />;
          },
        },
        {
          accessorKey: "reset_time",
          header: t("resetTime"),
          cell: ({ row }) => {
            return <Display type="number" value={row.getValue("reset_time")} unlimited />;
          },
        },
        {
          accessorKey: "expire_time",
          header: t("expireTime"),
          cell: ({ row }) =>
            row.getValue("expire_time") ? formatDate(row.getValue("expire_time")) : t("permanent"),
        },
        {
          accessorKey: "created_at",
          header: t("createdAt"),
          cell: ({ row }) => formatDate(row.getValue("created_at")),
        },
      ]}
      request={async (pagination) => {
        const { data } = await getUserSubscribe({
          query: {
            user_id: userId,
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
          return [
            <SubscriptionForm
              key="edit"
              trigger={t("edit")}
              title={t("editSubscription")}
              loading={loading}
              initialData={row}
              onSubmit={async (values) => {
                setLoading(true);
                await updateUserSubscribe({
                  body: {
                    user_subscribe_id: row.id,
                    subscribe_id: values.subscribe_id ?? 0,
                    traffic: values.traffic ?? 0,
                    expired_at: values.expired_at ?? 0,
                    upload: values.upload ?? 0,
                    download: values.download ?? 0,
                  },
                });
                toast.success(t("updateSuccess"));
                ref.current?.refresh();
                setLoading(false);
                return true;
              }}
            />,
            <Button
              key="copy"
              variant="secondary"
              onClick={async () => {
                await navigator.clipboard.writeText(getUserSubscribeUrls(row.token)[0] || "");
                toast.success(t("copySuccess"));
              }}
            >
              {t("copySubscription")}
            </Button>,
            <ConfirmButton
              key="delete"
              trigger={<Button variant="destructive">{t("delete")}</Button>}
              title={t("confirmDelete")}
              description={t("deleteSubscriptionDescription")}
              onConfirm={async () => {
                await deleteUserSubscribe({ body: { user_subscribe_id: row.id } });
                toast.success(t("deleteSuccess"));
                ref.current?.refresh();
              }}
              cancelText={t("cancel")}
              confirmText={t("confirm")}
            />,
            <RowMoreActions key="more" userId={userId} subId={row.id} />,
          ];
        },
      }}
    />
  );
}

function RowMoreActions({ userId, subId }: { userId: number; subId: number }) {
  const triggerRef = useRef<HTMLButtonElement>(null);
  const t = useTranslations("user");
  return (
    <div className="inline-flex">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline">{t("more")}</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem asChild>
            <AdminLink
              href={`/dashboard/log/subscribe?user_id=${userId}&user_subscribe_id=${subId}`}
            >
              {t("subscriptionLogs")}
            </AdminLink>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <AdminLink
              href={`/dashboard/log/reset-subscribe?user_id=${userId}&user_subscribe_id=${subId}`}
            >
              {t("resetLogs")}
            </AdminLink>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <AdminLink
              href={`/dashboard/log/subscribe-traffic?user_id=${userId}&user_subscribe_id=${subId}`}
            >
              {t("trafficStats")}
            </AdminLink>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <AdminLink
              href={`/dashboard/log/traffic-details?user_id=${userId}&subscribe_id=${subId}`}
            >
              {t("trafficDetails")}
            </AdminLink>
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={(e) => {
              e.preventDefault();
              triggerRef.current?.click();
            }}
          >
            {t("onlineDevices")}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <SubscriptionDetail
        trigger={<Button ref={triggerRef} className="hidden" />}
        userId={userId}
        subscriptionId={subId}
      />
    </div>
  );
}
