'use client';

import { UserDetail, UserSubscribeDetail } from '@/app/dashboard/user/user-detail';
import { OrderLink } from '@/components/order-link';
import { ProTable } from '@/components/pro-table';
import { filterResetSubscribeLog } from '@/services/admin/log';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';

export default function ResetSubscribeLogPage() {
  const t = useTranslations('log');
  const sp = useSearchParams();

  const today = new Date().toISOString().split('T')[0];

  const getResetSubscribeTypeText = (type: number) => {
    const typeText = t(`type.${type}`);
    if (typeText === `log.type.${type}`) {
      return `${t('unknown')} (${type})`;
    }
    return typeText;
  };

  const initialFilters = {
    date: sp.get('date') || today,
    user_subscribe_id: sp.get('user_subscribe_id')
      ? Number(sp.get('user_subscribe_id'))
      : undefined,
  };
  return (
    <ProTable<API.ResetSubscribeLog, { date?: string; user_subscribe_id?: number }>
      header={{ title: t('title.resetSubscribe') }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: 'user',
          header: t('column.user'),
          cell: ({ row }) => <UserDetail id={Number(row.original.user_id)} />,
        },
        {
          accessorKey: 'user_subscribe_id',
          header: t('column.subscribeId'),
          cell: ({ row }) => (
            <UserSubscribeDetail id={Number(row.original.user_subscribe_id)} enabled hoverCard />
          ),
        },
        {
          accessorKey: 'type',
          header: t('column.type'),
          cell: ({ row }) => <Badge>{getResetSubscribeTypeText(row.original.type)}</Badge>,
        },
        {
          accessorKey: 'order_no',
          header: t('column.orderNo'),
          cell: ({ row }) => <OrderLink orderId={row.original.order_no} />,
        },
        {
          accessorKey: 'timestamp',
          header: t('column.time'),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: 'date', type: 'date' },
        { key: 'user_subscribe_id', placeholder: t('column.subscribeId') },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterResetSubscribeLog({
          page: pagination.page,
          size: pagination.size,
          date: (filter as any)?.date,
          user_subscribe_id: (filter as any)?.user_subscribe_id,
        });
        const list = (data?.data?.list || []) as any[];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
