'use client';

import { UserDetail, UserSubscribeDetail } from '@/app/dashboard/user/user-detail';
import { ProTable } from '@/components/pro-table';
import { filterTrafficLogDetails } from '@/services/admin/log';
import { useServer } from '@/store/server';
import { formatDate } from '@/utils/common';
import { formatBytes } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';

export default function TrafficDetailsPage() {
  const t = useTranslations('log');
  const sp = useSearchParams();
  const { getServerName } = useServer();

  const today = new Date().toISOString().split('T')[0];

  const initialFilters = {
    date: sp.get('date') || today,
    server_id: sp.get('server_id') ? Number(sp.get('server_id')) : undefined,
    user_id: sp.get('user_id') ? Number(sp.get('user_id')) : undefined,
    subscribe_id: sp.get('subscribe_id') ? Number(sp.get('subscribe_id')) : undefined,
  };
  return (
    <ProTable<API.TrafficLogDetails, { search?: string }>
      header={{ title: t('title.trafficDetails') }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: 'server_id',
          header: t('column.server'),
          cell: ({ row }) => (
            <span>
              {getServerName(row.original.server_id)} ({row.original.server_id})
            </span>
          ),
        },
        {
          accessorKey: 'user_id',
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
          accessorKey: 'timestamp',
          header: t('column.time'),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: 'date', type: 'date' },
        { key: 'server_id', placeholder: t('column.serverId') },
        { key: 'user_id', placeholder: t('column.userId') },
        { key: 'subscribe_id', placeholder: t('column.subscribeId') },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterTrafficLogDetails({
          page: pagination.page,
          size: pagination.size,
          date: (filter as any)?.date,
          server_id: (filter as any)?.server_id,
          user_id: (filter as any)?.user_id,
          subscribe_id: (filter as any)?.subscribe_id,
        });
        const list = (data?.data?.list || []) as any[];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
