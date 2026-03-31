'use client';

import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createServer,
  deleteServer,
  filterServerList,
  resetSortWithServer,
  updateServer,
} from '@/services/admin/server';
import { useNode } from '@/store/node';
import { useServer } from '@/store/server';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { cn } from '@workspace/ui/lib/utils';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import DynamicMultiplier from './dynamic-multiplier';
import OnlineUsersCell from './online-users-cell';
import ServerConfig from './server-config';
import ServerForm from './server-form';
import ServerInstall from './server-install';

function PctBar({ value }: { value: number }) {
  const v = value.toFixed(2);
  return (
    <div className='min-w-24'>
      <div className='text-xs leading-none'>{v}%</div>
      <div className='bg-muted h-1.5 w-full rounded'>
        <div className='bg-primary h-1.5 rounded' style={{ width: `${v}%` }} />
      </div>
    </div>
  );
}

function RegionIpCell({
  country,
  city,
  ip,
  t,
}: {
  country?: string;
  city?: string;
  ip?: string;
  t: (key: string) => string;
}) {
  const region = [country, city].filter(Boolean).join(' / ') || t('notAvailable');
  return (
    <div className='flex items-center gap-1'>
      <Badge variant='outline'>{region}</Badge>
      <Badge variant='secondary'>{ip || t('notAvailable')}</Badge>
    </div>
  );
}

export default function ServersPage() {
  const t = useTranslations('servers');
  const { isServerReferencedByNodes } = useNode();
  const { fetchServers } = useServer();

  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);

  return (
    <div className='space-y-4'>
      <div className='grid grid-cols-1 gap-4 md:grid-cols-2'>
        <DynamicMultiplier />
        <ServerConfig />
      </div>
      <ProTable<API.Server, { search: string }>
        action={ref}
        header={{
          title: t('pageTitle'),
          toolbar: (
            <div className='flex gap-2'>
              <ServerForm
                trigger={t('create')}
                title={t('drawerCreateTitle')}
                loading={loading}
                onSubmit={async (values) => {
                  setLoading(true);
                  try {
                    await createServer(values as unknown as API.CreateServerRequest);
                    toast.success(t('created'));
                    ref.current?.refresh();
                    fetchServers();
                    setLoading(false);
                    return true;
                  } catch (e) {
                    setLoading(false);
                    return false;
                  }
                }}
              />
            </div>
          ),
        }}
        columns={[
          {
            accessorKey: 'id',
            header: t('id'),
            cell: ({ row }) => <Badge>{row.getValue('id')}</Badge>,
          },
          { accessorKey: 'name', header: t('name') },
          {
            id: 'region_ip',
            header: t('address'),
            cell: ({ row }) => (
              <RegionIpCell
                country={row.original.country as unknown as string}
                city={row.original.city as unknown as string}
                ip={row.original.address as unknown as string}
                t={t}
              />
            ),
          },
          {
            accessorKey: 'protocols',
            header: t('protocols'),
            cell: ({ row }) => {
              const list = row.original.protocols.filter((p) => p.enable) as API.Protocol[];
              if (!list.length) return '—';
              return (
                <div className='flex flex-col gap-1'>
                  {list.map((p, idx) => {
                    const ratio = Number(p.ratio ?? 1) || 1;
                    return (
                      <div key={idx} className='flex items-center gap-2'>
                        <Badge variant='outline'>{ratio.toFixed(2)}x</Badge>
                        <Badge variant='secondary'>{p.type}</Badge>
                        <Badge variant='secondary'>{p.port}</Badge>
                      </div>
                    );
                  })}
                </div>
              );
            },
          },

          {
            id: 'status',
            header: t('status'),
            cell: ({ row }) => {
              const offline = row.original.status.status === 'offline';
              return (
                <div className='flex items-center gap-2'>
                  <span
                    className={cn(
                      'inline-block h-2.5 w-2.5 rounded-full',
                      offline ? 'bg-zinc-400' : 'bg-emerald-500',
                    )}
                  />
                  <span className='text-sm'>{offline ? t('offline') : t('online')}</span>
                </div>
              );
            },
          },
          {
            id: 'cpu',
            header: t('cpu'),
            cell: ({ row }) => (
              <PctBar value={(row.original.status?.cpu as unknown as number) ?? 0} />
            ),
          },
          {
            id: 'mem',
            header: t('memory'),
            cell: ({ row }) => (
              <PctBar value={(row.original.status?.mem as unknown as number) ?? 0} />
            ),
          },
          {
            id: 'disk',
            header: t('disk'),
            cell: ({ row }) => (
              <PctBar value={(row.original.status?.disk as unknown as number) ?? 0} />
            ),
          },

          {
            id: 'online_users',
            header: t('onlineUsers'),
            cell: ({ row }) => <OnlineUsersCell status={row.original.status as API.ServerStatus} />,
          },
          // traffic ratio moved to per-protocol configs; column removed
        ]}
        params={[{ key: 'search' }]}
        request={async (pagination, filter) => {
          const { data } = await filterServerList({
            page: pagination.page,
            size: pagination.size,
            search: filter?.search || undefined,
          });
          const list = (data?.data?.list || []) as API.Server[];
          const total = (data?.data?.total ?? list.length) as number;
          return { list, total };
        }}
        actions={{
          render: (row) => [
            <ServerForm
              key='edit'
              trigger={t('edit')}
              title={t('drawerEditTitle')}
              initialValues={row}
              loading={loading}
              onSubmit={async (values) => {
                setLoading(true);
                try {
                  // ServerForm already returns API-shaped body; add id for update
                  await updateServer({
                    id: row.id,
                    ...(values as unknown as Omit<API.UpdateServerRequest, 'id'>),
                  });
                  toast.success(t('updated'));
                  ref.current?.refresh();
                  fetchServers();
                  setLoading(false);
                  return true;
                } catch (e) {
                  setLoading(false);
                  return false;
                }
              }}
            />,
            <ServerInstall key='install' server={row} />,
            <ConfirmButton
              key='delete'
              trigger={
                <Button variant='destructive' disabled={isServerReferencedByNodes(row.id)}>
                  {t('delete')}
                </Button>
              }
              title={t('confirmDeleteTitle')}
              description={t('confirmDeleteDesc')}
              onConfirm={async () => {
                await deleteServer({ id: row.id } as any);
                toast.success(t('deleted'));
                ref.current?.refresh();
                fetchServers();
              }}
              cancelText={t('cancel')}
              confirmText={t('confirm')}
            />,
            <Button
              key='copy'
              variant='outline'
              onClick={async () => {
                setLoading(true);
                const { id, created_at, updated_at, last_reported_at, status, ...others } =
                  row as any;
                const body: API.CreateServerRequest = {
                  name: others.name,
                  country: others.country,
                  city: others.city,
                  address: others.address,
                  protocols: others.protocols || [],
                };
                await createServer(body);
                toast.success(t('copied'));
                ref.current?.refresh();
                fetchServers();
                setLoading(false);
              }}
            >
              {t('copy')}
            </Button>,
          ],
          batchRender(rows) {
            const hasReferencedServers = rows.some((row) => isServerReferencedByNodes(row.id));
            return [
              <ConfirmButton
                key='delete'
                trigger={
                  <Button variant='destructive' disabled={hasReferencedServers}>
                    {t('delete')}
                  </Button>
                }
                title={t('confirmDeleteTitle')}
                description={t('confirmDeleteDesc')}
                onConfirm={async () => {
                  await Promise.all(rows.map((r) => deleteServer({ id: r.id })));
                  toast.success(t('deleted'));
                  ref.current?.refresh();
                  fetchServers();
                }}
                cancelText={t('cancel')}
                confirmText={t('confirm')}
              />,
            ];
          },
        }}
        onSort={async (source, target, items) => {
          const sourceIndex = items.findIndex((item) => String(item.id) === source);
          const targetIndex = items.findIndex((item) => String(item.id) === target);

          const originalSorts = items.map((item) => item.sort);

          const [movedItem] = items.splice(sourceIndex, 1);
          items.splice(targetIndex, 0, movedItem!);

          const updatedItems = items.map((item, index) => {
            const originalSort = originalSorts[index];
            const newSort = originalSort !== undefined ? originalSort : item.sort;
            return { ...item, sort: newSort };
          });

          const changedItems = updatedItems.filter((item, index) => {
            return item.sort !== items[index]?.sort;
          });

          if (changedItems.length > 0) {
            resetSortWithServer({
              sort: changedItems.map((item) => ({
                id: item.id,
                sort: item.sort,
              })) as API.SortItem[],
            });
            toast.success(t('sorted_success'));
          }
          return updatedItems;
        }}
      />
    </div>
  );
}
