'use client';

import { UserDetail, UserSubscribeDetail } from '@/app/dashboard/user/user-detail';
import { ProTable } from '@/components/pro-table';
import { filterUserSubscribeTrafficLog } from '@/services/admin/log';
import { Button } from '@workspace/ui/components/button';
import { formatBytes } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';

export default function SubscribeTrafficLogPage() {
  const t = useTranslations('log');
  const sp = useSearchParams();

  const today = new Date().toISOString().split('T')[0];

  const initialFilters = {
    date: sp.get('date') || today,
    user_id: sp.get('user_id') ? Number(sp.get('user_id')) : undefined,
    user_subscribe_id: sp.get('user_subscribe_id')
      ? Number(sp.get('user_subscribe_id'))
      : undefined,
  };
  return (
    <ProTable<
      API.UserSubscribeTrafficLog,
      { date?: string; user_id?: number; user_subscribe_id?: number }
    >
      header={{ title: t('title.subscribeTraffic') }}
      initialFilters={initialFilters}
      actions={{
        render: (row) => [
          <Button key='detail' asChild>
            <Link
              href={`/dashboard/log/traffic-details?date=${row.date}&user_id=${row.user_id}&subscribe_id=${row.subscribe_id}`}
            >
              {t('detail')}
            </Link>
          </Button>,
        ],
      }}
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
          accessorKey: 'upload',
          header: t('column.upload'),
          cell: ({ row }) => formatBytes(row.original.upload),
        },
        {
          accessorKey: 'download',
          header: t('column.download'),
          cell: ({ row }) => formatBytes(row.original.download),
        },
        {
          accessorKey: 'total',
          header: t('column.total'),
          cell: ({ row }) => formatBytes(row.original.total),
        },
        {
          accessorKey: 'date',
          header: t('column.date'),
        },
      ]}
      params={[
        { key: 'date', type: 'date' },
        { key: 'user_id', placeholder: t('column.userId') },
        { key: 'user_subscribe_id', placeholder: t('column.subscribeId') },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterUserSubscribeTrafficLog({
          page: pagination.page,
          size: pagination.size,
          date: (filter as any)?.date,
          user_id: (filter as any)?.user_id,
          user_subscribe_id: (filter as any)?.user_subscribe_id,
        });
        const list = ((data?.data?.list || []) as API.UserSubscribeTrafficLog[]) || [];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
