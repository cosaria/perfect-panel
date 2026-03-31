'use client';

import { Display } from '@/components/display';
import { ProTable } from '@/components/pro-table';
import { queryQuotaTaskList } from '@/services/admin/marketing';
import { useSubscribe } from '@/store/subscribe';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useState } from 'react';

export default function QuotaTaskManager() {
  const t = useTranslations('marketing');
  const [open, setOpen] = useState(false);

  const { subscribes } = useSubscribe();
  const subscribeMap =
    subscribes?.reduce(
      (acc, subscribe) => {
        acc[subscribe.id!] = subscribe.name!;
        return acc;
      },
      {} as Record<number, string>,
    ) || {};

  const getStatusBadge = (status: number) => {
    const statusConfig = {
      0: { label: t('notStarted'), variant: 'secondary' as const },
      1: { label: t('inProgress'), variant: 'default' as const },
      2: { label: t('completed'), variant: 'default' as const },
    };

    const config = statusConfig[status as keyof typeof statusConfig] || {
      label: `${t('status')} ${status}`,
      variant: 'secondary' as const,
    };

    return <Badge variant={config.variant}>{config.label}</Badge>;
  };

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <div className='flex cursor-pointer items-center justify-between transition-colors'>
          <div className='flex items-center gap-3'>
            <div className='bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg'>
              <Icon icon='mdi:database-plus' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('quotaTaskManager')}</p>
              <p className='text-muted-foreground text-sm'>{t('viewAndManageQuotaTasks')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[1000px] max-w-full md:max-w-screen-lg'>
        <SheetHeader>
          <SheetTitle>{t('quotaTasks')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-env(safe-area-inset-top))] px-6'>
          <div className='mt-4 space-y-4'>
            {open && (
              <ProTable<API.QuotaTask, API.QueryQuotaTaskListParams>
                columns={[
                  {
                    accessorKey: 'subscribers',
                    header: t('subscribers'),
                    size: 200,
                    cell: ({ row }) => {
                      const subscribers = row.getValue('subscribers') as number[];
                      const subscriptionNames =
                        subscribers?.map((id) => subscribeMap[id]).filter(Boolean) || [];

                      if (subscriptionNames.length === 0) {
                        return (
                          <span className='text-muted-foreground text-sm'>
                            {t('noSubscriptions')}
                          </span>
                        );
                      }

                      return (
                        <div className='flex flex-wrap gap-1'>
                          {subscriptionNames.map((name, index) => (
                            <span key={index} className='bg-muted rounded px-2 py-1 text-xs'>
                              {name}
                            </span>
                          ))}
                        </div>
                      );
                    },
                  },
                  {
                    accessorKey: 'is_active',
                    header: t('validOnly'),
                    size: 120,
                    cell: ({ row }) => {
                      const isActive = row.getValue('is_active') as boolean;
                      return <span className='text-sm'>{isActive ? t('yes') : t('no')}</span>;
                    },
                  },
                  {
                    accessorKey: 'reset_traffic',
                    header: t('resetTraffic'),
                    size: 120,
                    cell: ({ row }) => {
                      const resetTraffic = row.getValue('reset_traffic') as boolean;
                      return <span className='text-sm'>{resetTraffic ? t('yes') : t('no')}</span>;
                    },
                  },
                  {
                    accessorKey: 'gift_value',
                    header: t('giftAmount'),
                    size: 120,
                    cell: ({ row }) => {
                      const giftValue = row.getValue('gift_value') as number;
                      const task = row.original as API.QuotaTask;
                      const giftType = task.gift_type;

                      return (
                        <div className='text-sm font-medium'>
                          {giftType === 1 ? (
                            <Display type='currency' value={giftValue} />
                          ) : (
                            `${giftValue}%`
                          )}
                        </div>
                      );
                    },
                  },
                  {
                    accessorKey: 'days',
                    header: t('quotaDays'),
                    size: 100,
                    cell: ({ row }) => {
                      const days = row.getValue('days') as number;
                      return (
                        <span className='font-medium'>
                          {days} {t('days')}
                        </span>
                      );
                    },
                  },
                  {
                    accessorKey: 'time_range',
                    header: t('timeRange'),
                    size: 180,
                    cell: ({ row }) => {
                      const task = row.original as API.QuotaTask;
                      const startTime = task.start_time;
                      const endTime = task.end_time;

                      if (!startTime && !endTime) {
                        return (
                          <span className='text-muted-foreground text-sm'>{t('noTimeLimit')}</span>
                        );
                      }

                      return (
                        <div className='space-y-1 text-xs'>
                          {startTime && (
                            <div>
                              {t('startTime')}: {formatDate(startTime)}
                            </div>
                          )}
                          {endTime && (
                            <div>
                              {t('endTime')}: {formatDate(endTime)}
                            </div>
                          )}
                        </div>
                      );
                    },
                  },
                  {
                    accessorKey: 'status',
                    header: t('status'),
                    size: 100,
                    cell: ({ row }) => getStatusBadge(row.getValue('status') as number),
                  },
                  {
                    accessorKey: 'created_at',
                    header: t('createdAt'),
                    size: 150,
                    cell: ({ row }) => {
                      const createdAt = row.getValue('created_at') as number;
                      return formatDate(createdAt);
                    },
                  },
                ]}
                request={async (pagination, filters) => {
                  const response = await queryQuotaTaskList({
                    ...filters,
                    page: pagination.page,
                    size: pagination.size,
                  });
                  return {
                    list: response.data?.data?.list || [],
                    total: response.data?.data?.total || 0,
                  };
                }}
                params={[
                  {
                    key: 'status',
                    placeholder: t('status'),
                    options: [
                      { label: t('notStarted'), value: '0' },
                      { label: t('inProgress'), value: '1' },
                      { label: t('completed'), value: '2' },
                    ],
                  },
                ]}
              />
            )}
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
