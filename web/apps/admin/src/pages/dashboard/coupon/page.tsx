"use client";

import { Badge } from "@workspace/ui/components/badge";
import { Button } from "@workspace/ui/components/button";
import { Switch } from "@workspace/ui/components/switch";
import { ConfirmButton } from "@workspace/ui/custom-components/confirm-button";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { useRef, useState } from "react";
import { toast } from "sonner";
import { Display } from "@/components/display";
import { ProTable, type ProTableActions } from "@/components/pro-table";
import {
  batchDeleteCoupon,
  createCoupon,
  deleteCoupon,
  getCouponList,
  updateCoupon,
} from "@/services/admin-api/sdk.gen";
import type {
  Coupon,
  CreateCouponRequest,
  UpdateCouponRequest,
} from "@/services/admin-api/types.gen";
import { useSubscribe } from "@/store/subscribe";
import { formatDate } from "@/utils/common";
import CouponForm from "./coupon-form";

export default function Page() {
  const t = useTranslations("coupon");
  const [loading, setLoading] = useState(false);
  const { subscribes } = useSubscribe();
  const ref = useRef<ProTableActions>(null);
  return (
    <ProTable<Coupon, { group_id: number; query: string }>
      action={ref}
      header={{
        toolbar: (
          <CouponForm<CreateCouponRequest>
            trigger={t("create")}
            title={t("createCoupon")}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createCoupon({
                  body: {
                    ...values,
                    enable: false,
                  },
                });
                toast.success(t("createSuccess"));
                ref.current?.refresh();
                setLoading(false);
                return true;
              } catch (_error) {
                setLoading(false);
                return false;
              }
            }}
          />
        ),
      }}
      params={[
        {
          key: "subscribe",
          placeholder: t("subscribe"),
          options: subscribes
            ?.filter(
              (item): item is typeof item & { id: number; name: string } =>
                typeof item.id === "number" && typeof item.name === "string",
            )
            .map((item) => ({
              label: item.name,
              value: String(item.id),
            })),
        },
        {
          key: "search",
        },
      ]}
      request={async (pagination, filters) => {
        const { data } = await getCouponList({
          query: {
            ...pagination,
            ...filters,
          },
        });
        return {
          list: data?.list || [],
          total: data?.total || 0,
        };
      }}
      columns={[
        {
          accessorKey: "enable",
          header: t("enable"),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue("enable")}
                onCheckedChange={async (checked) => {
                  await updateCoupon({
                    body: {
                      ...row.original,
                      enable: checked,
                    } as UpdateCouponRequest,
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: "name",
          header: t("name"),
        },
        {
          accessorKey: "code",
          header: t("code"),
        },
        {
          accessorKey: "type",
          header: t("type"),
          cell: ({ row }) => (
            <Badge variant={row.getValue("type") === 1 ? "default" : "secondary"}>
              {row.getValue("type") === 1 ? t("percentage") : t("amount")}
            </Badge>
          ),
        },
        {
          accessorKey: "discount",
          header: t("discount"),
          cell: ({ row }) => (
            <Badge variant={row.getValue("type") === 1 ? "default" : "secondary"}>
              {row.getValue("type") === 1 ? (
                `${row.original.discount} %`
              ) : (
                <Display type="currency" value={row.original.discount} />
              )}
            </Badge>
          ),
        },
        {
          accessorKey: "count",
          header: t("count"),
          cell: ({ row }) => (
            <div className="flex flex-col">
              <span>
                {t("count")}: {row.original.count === 0 ? t("unlimited") : row.original.count}
              </span>
              <span>
                {t("remainingTimes")}:{" "}
                {row.original.count === 0
                  ? t("unlimited")
                  : row.original.count - row.original.used_count}
              </span>
              <span>
                {t("usedTimes")}: {row.original.used_count}
              </span>
            </div>
          ),
        },
        {
          accessorKey: "expire",
          header: t("validityPeriod"),
          cell: ({ row }) => {
            const { start_time, expire_time } = row.original;
            if (start_time) {
              return expire_time ? (
                <>
                  {formatDate(start_time)} - {formatDate(expire_time)}
                </>
              ) : start_time ? (
                formatDate(start_time)
              ) : (
                "--"
              );
            }
            return "--";
          },
        },
      ]}
      actions={{
        render: (row) => [
          <CouponForm<UpdateCouponRequest>
            key="edit"
            trigger={t("edit")}
            title={t("editCoupon")}
            loading={loading}
            initialValues={row}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await updateCoupon({ body: { ...row, ...values } });
                toast.success(t("updateSuccess"));
                ref.current?.refresh();
                setLoading(false);
                return true;
              } catch (_error) {
                setLoading(false);
                return false;
              }
            }}
          />,
          <ConfirmButton
            key="delete"
            trigger={<Button variant="destructive">{t("delete")}</Button>}
            title={t("confirmDelete")}
            description={t("deleteWarning")}
            onConfirm={async () => {
              await deleteCoupon({ body: { id: row.id } });
              toast.success(t("deleteSuccess"));
              ref.current?.refresh();
            }}
            cancelText={t("cancel")}
            confirmText={t("confirm")}
          />,
        ],
        batchRender: (rows) => [
          <ConfirmButton
            key="delete"
            trigger={<Button variant="destructive">{t("delete")}</Button>}
            title={t("confirmDelete")}
            description={t("deleteWarning")}
            onConfirm={async () => {
              await batchDeleteCoupon({ body: { ids: rows.map((item) => item.id) } });
              toast.success(t("deleteSuccess"));
              ref.current?.reset();
            }}
            cancelText={t("cancel")}
            confirmText={t("confirm")}
          />,
        ],
      }}
    />
  );
}
