'use client';

import { UserDetail, UserSubscribeDetail } from '@/app/dashboard/user/user-detail';
import { Display } from '@/components/display';
import { OrderLink } from '@/components/order-link';
import { ProTable } from '@/components/pro-table';
import { filterGiftLog } from '@/services/admin/log';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';

export default function GiftLogPage() {
  const t = useTranslations('log');
  const sp = useSearchParams();

  const today = new Date().toISOString().split('T')[0];

  const getGiftTypeText = (type: number) => {
    const typeText = t(`type.${type}`);
    if (typeText === `log.type.${type}`) {
      return `${t('unknown')} (${type})`;
    }
    return typeText;
  };

  const initialFilters = {
    date: sp.get('date') || today,
    user_id: sp.get('user_id') ? Number(sp.get('user_id')) : undefined,
  };
  return (
    <ProTable<API.GiftLog, { search?: string }>
      header={{ title: t('title.gift') }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: 'user',
          header: t('column.user'),
          cell: ({ row }) => <UserDetail id={Number(row.original.user_id)} />,
        },
        {
          accessorKey: 'subscribe_id',
          header: t('column.subscribe'),
          cell: ({ row }) => (
            <UserSubscribeDetail id={Number(row.original.subscribe_id)} enabled hoverCard />
          ),
        },
        {
          accessorKey: 'order_no',
          header: t('column.orderNo'),
          cell: ({ row }) => <OrderLink orderId={row.original.order_no} />,
        },
        {
          accessorKey: 'amount',
          header: t('column.amount'),
          cell: ({ row }) => <Display type='currency' value={row.original.amount} />,
        },
        {
          accessorKey: 'balance',
          header: t('column.balance'),
          cell: ({ row }) => <Display type='currency' value={row.original.balance} />,
        },
        {
          accessorKey: 'type',
          header: t('column.type'),
          cell: ({ row }) => <Badge>{getGiftTypeText(row.original.type)}</Badge>,
        },
        { accessorKey: 'remark', header: t('column.remark') },
        {
          accessorKey: 'timestamp',
          header: t('column.time'),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: 'date', type: 'date' },
        { key: 'user_id', placeholder: t('column.userId') },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterGiftLog({
          page: pagination.page,
          size: pagination.size,
          date: (filter as any)?.date,
          user_id: (filter as any)?.user_id,
        });
        const list = (data?.data?.list || []) as any[];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
