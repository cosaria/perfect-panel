'use client';

import { UserDetail, UserSubscribeDetail } from '@/app/dashboard/user/user-detail';
import { IpLink } from '@/components/ip-link';
import { ProTable } from '@/components/pro-table';
import { filterSubscribeLog } from '@/services/admin/log';
import { formatDate } from '@/utils/common';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@workspace/ui/components/tooltip';
import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';

export default function SubscribeLogPage() {
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
    <ProTable<API.SubscribeLog, { date?: string; user_id?: number }>
      header={{ title: t('title.subscribe') }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: 'user',
          header: t('column.user'),
          cell: ({ row }) => <UserDetail id={Number(row.original.user_id)} />,
        },
        {
          accessorKey: 'user_subscribe_id',
          header: t('column.subscribe'),
          cell: ({ row }) => (
            <UserSubscribeDetail id={Number(row.original.user_subscribe_id)} enabled hoverCard />
          ),
        },
        {
          accessorKey: 'client_ip',
          header: t('column.ip'),
          cell: ({ row }) => <IpLink ip={String((row.original as any).client_ip || '')} />,
        },
        {
          accessorKey: 'user_agent',
          header: t('column.userAgent'),
          cell: ({ row }) => {
            const userAgent = String(row.original.user_agent || '');
            return (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div className='max-w-48 cursor-help truncate'>{userAgent}</div>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p className='max-w-md break-words'>{userAgent}</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            );
          },
        },
        {
          accessorKey: 'timestamp',
          header: t('column.time'),
          cell: ({ row }) => formatDate(row.original.timestamp),
        },
      ]}
      params={[
        { key: 'date', type: 'date' },
        { key: 'user_id', placeholder: t('column.userId') },
        { key: 'user_subscribe_id', placeholder: t('column.subscribeId') },
      ]}
      request={async (pagination, filter) => {
        const { data } = await filterSubscribeLog({
          page: pagination.page,
          size: pagination.size,
          date: (filter as any)?.date,
          user_id: (filter as any)?.user_id,
          user_subscribe_id: (filter as any)?.user_subscribe_id,
        });
        const list = (data?.data?.list || []) as any[];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
