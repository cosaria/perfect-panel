'use client';

import { UserDetail } from '@/app/dashboard/user/user-detail';
import { IpLink } from '@/components/ip-link';
import { ProTable } from '@/components/pro-table';
import { getUserSubscribeById } from '@/services/admin/user';
import { formatDate } from '@/utils/common';
import { useQuery } from '@tanstack/react-query';
import { Badge } from '@workspace/ui/components/badge';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { formatBytes } from '@workspace/ui/utils';
import { Users } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useState } from 'react';

function UserSubscribeInfo({
  subscribeId,
  open,
  type,
}: {
  subscribeId: number;
  open: boolean;
  type: 'account' | 'subscribeName' | 'subscribeId' | 'trafficUsage' | 'expireTime';
}) {
  const t = useTranslations('servers');
  const { data } = useQuery({
    enabled: subscribeId !== 0 && open,
    queryKey: ['getUserSubscribeById', subscribeId],
    queryFn: async () => {
      const { data } = await getUserSubscribeById({ id: subscribeId });
      return data.data;
    },
  });

  if (!data) return <span className='text-muted-foreground'>--</span>;

  switch (type) {
    case 'account':
      if (!data.user_id) return <span className='text-muted-foreground'>--</span>;
      return <UserDetail id={data.user_id} />;

    case 'subscribeName':
      if (!data.subscribe?.name) return <span className='text-muted-foreground'>--</span>;
      return <span className='text-sm'>{data.subscribe.name}</span>;

    case 'subscribeId':
      if (!data.id) return <span className='text-muted-foreground'>--</span>;
      return <span className='font-mono text-sm'>{data.id}</span>;

    case 'trafficUsage': {
      const usedTraffic = data.upload + data.download;
      const totalTraffic = data.traffic || 0;
      return (
        <div className='min-w-0 text-sm'>
          <div className='break-words'>
            {formatBytes(usedTraffic)} / {totalTraffic > 0 ? formatBytes(totalTraffic) : '无限制'}
          </div>
        </div>
      );
    }

    case 'expireTime': {
      if (!data.expire_time) return <span className='text-muted-foreground'>--</span>;
      const isExpired = data.expire_time < Date.now() / 1000;
      return (
        <div className='flex flex-col gap-1 sm:flex-row sm:items-center sm:gap-2'>
          <span className='text-sm'>{formatDate(data.expire_time)}</span>
          {isExpired && (
            <Badge variant='destructive' className='w-fit px-1 py-0 text-xs'>
              {t('expired')}
            </Badge>
          )}
        </div>
      );
    }

    default:
      return <span className='text-muted-foreground'>--</span>;
  }
}

export default function OnlineUsersCell({ status }: { status?: API.ServerStatus }) {
  const t = useTranslations('servers');
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <button className='hover:text-foreground text-muted-foreground flex items-center gap-2 bg-transparent p-0 text-sm'>
          <Users className='h-4 w-4' /> {status?.online.length}
        </button>
      </SheetTrigger>
      <SheetContent className='h-screen w-screen max-w-none sm:h-auto sm:w-[900px] sm:max-w-[90vw]'>
        <SheetHeader>
          <SheetTitle>{t('onlineUsers')}</SheetTitle>
        </SheetHeader>
        <div className='-mx-6 h-[calc(100vh-48px-16px)] overflow-y-auto px-6 py-4 sm:h-[calc(100dvh-48px-16px-env(safe-area-inset-top))]'>
          <ProTable<API.ServerOnlineUser, Record<string, unknown>>
            header={{ hidden: true }}
            columns={[
              {
                accessorKey: 'ip',
                header: t('ipAddresses'),
                cell: ({ row }) => {
                  const ips = row.original.ip;
                  return (
                    <div className='flex min-w-0 flex-col gap-1'>
                      {ips.map((item, i) => (
                        <div className='whitespace-nowrap text-sm' key={i}>
                          <Badge>{item.protocol}</Badge>
                          <IpLink ip={item.ip} className='ml-1 font-medium' />
                        </div>
                      ))}
                    </div>
                  );
                },
              },
              {
                accessorKey: 'user',
                header: t('user'),
                cell: ({ row }) => (
                  <UserSubscribeInfo
                    subscribeId={Number(row.original.subscribe_id)}
                    open={open}
                    type='account'
                  />
                ),
              },
              {
                accessorKey: 'subscription',
                header: t('subscription'),
                cell: ({ row }) => (
                  <UserSubscribeInfo
                    subscribeId={Number(row.original.subscribe_id)}
                    open={open}
                    type='subscribeName'
                  />
                ),
              },
              {
                accessorKey: 'subscribeId',
                header: t('subscribeId'),
                cell: ({ row }) => (
                  <UserSubscribeInfo
                    subscribeId={Number(row.original.subscribe_id)}
                    open={open}
                    type='subscribeId'
                  />
                ),
              },
              {
                accessorKey: 'traffic',
                header: t('traffic'),
                cell: ({ row }) => (
                  <UserSubscribeInfo
                    subscribeId={Number(row.original.subscribe_id)}
                    open={open}
                    type='trafficUsage'
                  />
                ),
              },
              {
                accessorKey: 'expireTime',
                header: t('expireTime'),
                cell: ({ row }) => (
                  <UserSubscribeInfo
                    subscribeId={Number(row.original.subscribe_id)}
                    open={open}
                    type='expireTime'
                  />
                ),
              },
            ]}
            request={async () => ({
              list: status?.online || [],
              total: status?.online?.length || 0,
            })}
          />
        </div>
      </SheetContent>
    </Sheet>
  );
}
