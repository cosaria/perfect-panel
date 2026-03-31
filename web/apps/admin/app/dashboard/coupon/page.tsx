'use client';

import { Display } from '@/components/display';
import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  batchDeleteCoupon,
  createCoupon,
  deleteCoupon,
  getCouponList,
  updateCoupon,
} from '@/services/admin/coupon';
import { useSubscribe } from '@/store/subscribe';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Switch } from '@workspace/ui/components/switch';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import CouponForm from './coupon-form';

export default function Page() {
  const t = useTranslations('coupon');
  const [loading, setLoading] = useState(false);
  const { subscribes } = useSubscribe();
  const ref = useRef<ProTableActions>(null);
  return (
    <ProTable<API.Coupon, { group_id: number; query: string }>
      action={ref}
      header={{
        toolbar: (
          <CouponForm<API.CreateCouponRequest>
            trigger={t('create')}
            title={t('createCoupon')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createCoupon({
                  ...values,
                  enable: false,
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
          key: 'subscribe',
          placeholder: t('subscribe'),
          options: subscribes?.map((item) => ({
            label: item.name!,
            value: String(item.id),
          })),
        },
        {
          key: 'search',
        },
      ]}
      request={async (pagination, filters) => {
        const { data } = await getCouponList({
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
          accessorKey: 'enable',
          header: t('enable'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('enable')}
                onCheckedChange={async (checked) => {
                  await updateCoupon({
                    ...row.original,
                    enable: checked,
                  } as API.UpdateCouponRequest);
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'name',
          header: t('name'),
        },
        {
          accessorKey: 'code',
          header: t('code'),
        },
        {
          accessorKey: 'type',
          header: t('type'),
          cell: ({ row }) => (
            <Badge variant={row.getValue('type') === 1 ? 'default' : 'secondary'}>
              {row.getValue('type') === 1 ? t('percentage') : t('amount')}
            </Badge>
          ),
        },
        {
          accessorKey: 'discount',
          header: t('discount'),
          cell: ({ row }) => (
            <Badge variant={row.getValue('type') === 1 ? 'default' : 'secondary'}>
              {row.getValue('type') === 1 ? (
                `${row.original.discount} %`
              ) : (
                <Display type='currency' value={row.original.discount} />
              )}
            </Badge>
          ),
        },
        {
          accessorKey: 'count',
          header: t('count'),
          cell: ({ row }) => (
            <div className='flex flex-col'>
              <span>
                {t('count')}: {row.original.count === 0 ? t('unlimited') : row.original.count}
              </span>
              <span>
                {t('remainingTimes')}:{' '}
                {row.original.count === 0
                  ? t('unlimited')
                  : row.original.count - row.original.used_count}
              </span>
              <span>
                {t('usedTimes')}: {row.original.used_count}
              </span>
            </div>
          ),
        },
        {
          accessorKey: 'expire',
          header: t('validityPeriod'),
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
                '--'
              );
            }
            return '--';
          },
        },
      ]}
      actions={{
        render: (row) => [
          <CouponForm<API.UpdateCouponRequest>
            key='edit'
            trigger={t('edit')}
            title={t('editCoupon')}
            loading={loading}
            initialValues={row}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await updateCoupon({ ...row, ...values });
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
              await deleteCoupon({ id: row.id });
              toast.success(t('deleteSuccess'));
              ref.current?.refresh();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
        ],
        batchRender: (rows) => [
          <ConfirmButton
            key='delete'
            trigger={<Button variant='destructive'>{t('delete')}</Button>}
            title={t('confirmDelete')}
            description={t('deleteWarning')}
            onConfirm={async () => {
              await batchDeleteCoupon({ ids: rows.map((item) => item.id) });
              toast.success(t('deleteSuccess'));
              ref.current?.reset();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
        ],
      }}
    />
  );
}
