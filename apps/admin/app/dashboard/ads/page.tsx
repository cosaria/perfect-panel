"use client";

import { Badge } from "@workspace/ui/components/badge";
import { Button } from "@workspace/ui/components/button";
import { Switch } from "@workspace/ui/components/switch";
import { ConfirmButton } from "@workspace/ui/custom-components/confirm-button";
import { useTranslations } from "next-intl";
import { useRef, useState } from "react";
import { toast } from "sonner";
import { ProTable, type ProTableActions } from "@/components/pro-table";
import { createAds, deleteAds, getAdsList, updateAds } from "@/services/admin-api/sdk.gen";
import type { Ads } from "@/services/admin-api/types.gen";
import { formatDate } from "@/utils/common";
import AdsForm from "./ads-form";

export default function Page() {
  const t = useTranslations("ads");
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);

  return (
    <ProTable<Ads, Record<string, unknown>>
      action={ref}
      header={{
        toolbar: (
          <AdsForm
            trigger={t("create")}
            title={t("createAds")}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createAds({
                  body: { ...values, status: 0 },
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
          key: "status",
          placeholder: t("status"),
          options: [
            { label: t("enabled"), value: "1" },
            { label: t("disabled"), value: "0" },
          ],
        },
        {
          key: "search",
        },
      ]}
      request={async (pagination, filters) => {
        const { data } = await getAdsList({
          body: { ...pagination, ...filters },
        });
        return {
          list: data?.list || [],
          total: data?.total || 0,
        };
      }}
      columns={[
        {
          accessorKey: "status",
          header: t("status"),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue("status") === 1}
                onCheckedChange={async (checked) => {
                  await updateAds({
                    body: { ...row.original, status: checked ? 1 : 0 },
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: "title",
          header: t("title"),
        },
        {
          accessorKey: "type",
          header: t("type"),
          cell: ({ row }) => {
            const type = row.original.type;
            return <Badge>{type}</Badge>;
          },
        },
        {
          accessorKey: "target_url",
          header: t("targetUrl"),
        },
        {
          accessorKey: "description",
          header: t("form.description"),
        },
        {
          accessorKey: "period",
          header: t("validityPeriod"),
          cell: ({ row }) => {
            const { start_time, end_time } = row.original;
            return (
              <>
                {formatDate(start_time)} - {formatDate(end_time)}
              </>
            );
          },
        },
      ]}
      actions={{
        render: (row) => [
          <AdsForm
            key="edit"
            trigger={t("edit")}
            title={t("editAds")}
            loading={loading}
            initialValues={{ ...row, type: row.type as "image" | "video" }}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await updateAds({ body: { ...row, ...values } });
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
              await deleteAds({ body: { id: row.id } });
              toast.success(t("deleteSuccess"));
              ref.current?.refresh();
            }}
            cancelText={t("cancel")}
            confirmText={t("confirm")}
          />,
        ],
      }}
    />
  );
}
