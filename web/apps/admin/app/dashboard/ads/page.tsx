'use client';

import { ProTable, ProTableActions } from '@/components/pro-table';
import { createAds, deleteAds, getAdsList, updateAds } from '@/services/admin/ads';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Switch } from '@workspace/ui/components/switch';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import AdsForm from './ads-form';

export default function Page() {
  const t = useTranslations('ads');
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);

  return (
    <ProTable<API.Ads, Record<string, unknown>>
      action={ref}
      header={{
        toolbar: (
          <AdsForm<API.CreateAdsRequest>
            trigger={t('create')}
            title={t('createAds')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createAds({
                  ...values,
                  status: 0,
                });
                toast.success(t('createSuccess'));
                ref.current?.refresh();
                setLoading(false);
                return true;
              } catch (error) {
                setLoading(false);
                return false;
              }
            }}
          />
        ),
      }}
      params={[
        {
          key: 'status',
          placeholder: t('status'),
          options: [
            { label: t('enabled'), value: '1' },
            { label: t('disabled'), value: '0' },
          ],
        },
        {
          key: 'search',
        },
      ]}
      request={async (pagination, filters) => {
        const { data } = await getAdsList({
          ...pagination,
          ...filters,
        });
        return {
          list: data.data?.list || [],
          total: data.data?.total || 0,
        };
      }}
      columns={[
        {
          accessorKey: 'status',
          header: t('status'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('status') === 1}
                onCheckedChange={async (checked) => {
                  await updateAds({
                    ...row.original,
                    status: checked ? 1 : 0,
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'title',
          header: t('title'),
        },
        {
          accessorKey: 'type',
          header: t('type'),
          cell: ({ row }) => {
            const type = row.original.type;
            return <Badge>{type}</Badge>;
          },
        },
        {
          accessorKey: 'target_url',
          header: t('targetUrl'),
        },
        {
          accessorKey: 'description',
          header: t('form.description'),
        },
        {
          accessorKey: 'period',
          header: t('validityPeriod'),
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
          <AdsForm<API.UpdateAdsRequest>
            key='edit'
            trigger={t('edit')}
            title={t('editAds')}
            loading={loading}
            initialValues={row}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await updateAds({ ...row, ...values });
                toast.success(t('updateSuccess'));
                ref.current?.refresh();
                setLoading(false);
                return true;
              } catch (error) {
                setLoading(false);
                return false;
              }
            }}
          />,
          <ConfirmButton
            key='delete'
            trigger={<Button variant='destructive'>{t('delete')}</Button>}
            title={t('confirmDelete')}
            description={t('deleteWarning')}
            onConfirm={async () => {
              await deleteAds({ id: row.id });
              toast.success(t('deleteSuccess'));
              ref.current?.refresh();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
        ],
      }}
    />
  );
}
