'use client';

import { ProTable } from '@/components/pro-table';
import { filterMobileLog } from '@/services/admin/log';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';
export default function MobileLogPage() {
  const t = useTranslations('log');
  const sp = useSearchParams();

  const today = new Date().toISOString().split('T')[0];

  const initialFilters = {
    search: sp.get('search') || undefined,
    date: sp.get('date') || today,
  };
  return (
    <ProTable<API.MessageLog, { search?: string }>
      header={{ title: t('title.mobile') }}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: 'platform',
          header: t('column.platform'),
          cell: ({ row }) => <Badge>{row.getValue('platform')}</Badge>,
        },
        { accessorKey: 'to', header: t('column.to') },
        { accessorKey: 'subject', header: t('column.subject') },
        {
          accessorKey: 'content',
          header: t('column.content'),
          cell: ({ row }) => (
            <pre className='max-w-[480px] overflow-auto whitespace-pre-wrap break-words text-xs'>
              {JSON.stringify(row.original.content || {}, null, 2)}
            </pre>
          ),
        },
        {
          accessorKey: 'status',
          header: t('column.status'),
          cell: ({ row }) => {
            const status = row.original.status;
            const getStatusVariant = (status: any) => {
              if (status === 1) {
                return 'default';
              } else if (status === 0) {
                return 'destructive';
              }
              return 'outline';
            };

            const getStatusText = (status: any) => {
              if (status === 1) return t('sent');
              if (status === 0) return t('failed');
              return t('unknown');
            };

            return <Badge variant={getStatusVariant(status)}>{getStatusText(status)}</Badge>;
          },
        },
        {
          accessorKey: 'created_at',
          header: t('column.time'),
          cell: ({ row }) => formatDate(row.original.created_at),
        },
      ]}
      params={[{ key: 'search' }, { key: 'date', type: 'date' }]}
      request={async (pagination, filter) => {
        const { data } = await filterMobileLog({
          page: pagination.page,
          size: pagination.size,
          search: filter?.search,
          date: (filter as any)?.date,
        });
        const list = ((data?.data?.list || []) as API.MessageLog[]) || [];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
    />
  );
}
