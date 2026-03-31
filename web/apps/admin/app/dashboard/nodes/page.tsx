'use client';

import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createNode,
  deleteNode,
  filterNodeList,
  resetSortWithNode,
  toggleNodeStatus,
  updateNode,
} from '@/services/admin/server';
import { useNode } from '@/store/node';
import { useServer } from '@/store/server';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Switch } from '@workspace/ui/components/switch';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import NodeForm from './node-form';

export default function NodesPage() {
  const t = useTranslations('nodes');
  const ref = useRef<ProTableActions>(null);
  const [loading, setLoading] = useState(false);

  // Use our zustand store for server data
  const { getServerName, getServerAddress, getProtocolPort } = useServer();
  const { fetchNodes, fetchTags } = useNode();

  return (
    <ProTable<API.Node, { search: string }>
      action={ref}
      header={{
        title: t('pageTitle'),
        toolbar: (
          <NodeForm
            trigger={t('create')}
            title={t('drawerCreateTitle')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                const body: API.CreateNodeRequest = {
                  name: values.name,
                  server_id: Number(values.server_id!),
                  protocol: values.protocol,
                  address: values.address,
                  port: Number(values.port!),
                  tags: values.tags || [],
                  enabled: false,
                };
                await createNode(body);
                toast.success(t('created'));
                ref.current?.refresh();
                fetchNodes();
                fetchTags();
                setLoading(false);
                return true;
              } catch (e) {
                setLoading(false);
                return false;
              }
            }}
          />
        ),
      }}
      columns={[
        {
          id: 'enabled',
          header: t('enabled'),
          cell: ({ row }) => (
            <Switch
              checked={row.original.enabled}
              onCheckedChange={async (v) => {
                await toggleNodeStatus({ id: row.original.id, enable: v });
                toast.success(v ? t('enabled_on') : t('enabled_off'));
                ref.current?.refresh();
                fetchNodes();
                fetchTags();
              }}
            />
          ),
        },
        { accessorKey: 'name', header: t('name') },

        {
          id: 'address_port',
          header: `${t('address')}:${t('port')}`,
          cell: ({ row }) => `${row.original.address || '—'}:${row.original.port || '—'}`,
        },

        {
          id: 'server_id',
          header: t('server'),
          cell: ({ row }) =>
            `${getServerName(row.original.server_id)}:${getServerAddress(row.original.server_id)}`,
        },
        {
          id: 'protocol',
          header: ` ${t('protocol')}:${t('port')}`,
          cell: ({ row }) =>
            `${row.original.protocol}:${getProtocolPort(row.original.server_id, row.original.protocol)}`,
        },
        {
          accessorKey: 'tags',
          header: t('tags'),
          cell: ({ row }) => (
            <div className='flex flex-wrap gap-1'>
              {(row.original.tags || []).length === 0
                ? '—'
                : row.original.tags.map((tg) => (
                    <Badge key={tg} variant='outline'>
                      {tg}
                    </Badge>
                  ))}
            </div>
          ),
        },
      ]}
      params={[{ key: 'search' }]}
      request={async (pagination, filter) => {
        const { data } = await filterNodeList({
          page: pagination.page,
          size: pagination.size,
          search: filter?.search || undefined,
        });
        const list = (data?.data?.list || []) as API.Node[];
        const total = Number(data?.data?.total || list.length);
        return { list, total };
      }}
      actions={{
        render: (row) => [
          <NodeForm
            key='edit'
            trigger={t('edit')}
            title={t('drawerEditTitle')}
            loading={loading}
            initialValues={row}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                const body: API.UpdateNodeRequest = {
                  ...row,
                  ...values,
                } as any;
                await updateNode(body);
                toast.success(t('updated'));
                ref.current?.refresh();
                fetchNodes();
                fetchTags();
                setLoading(false);
                return true;
              } catch (e) {
                setLoading(false);
                return false;
              }
            }}
          />,
          <ConfirmButton
            key='delete'
            trigger={<Button variant='destructive'>{t('delete')}</Button>}
            title={t('confirmDeleteTitle')}
            description={t('confirmDeleteDesc')}
            onConfirm={async () => {
              await deleteNode({ id: row.id } as any);
              toast.success(t('deleted'));
              ref.current?.refresh();
              fetchNodes();
              fetchTags();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
          <Button
            key='copy'
            variant='outline'
            onClick={async () => {
              const { id, enabled, created_at, updated_at, sort, ...rest } = row as any;
              await createNode({
                ...rest,
                enabled: false,
              });
              toast.success(t('copied'));
              ref.current?.refresh();
              fetchNodes();
              fetchTags();
            }}
          >
            {t('copy')}
          </Button>,
        ],
        batchRender(rows) {
          return [
            <ConfirmButton
              key='delete'
              trigger={<Button variant='destructive'>{t('delete')}</Button>}
              title={t('confirmDeleteTitle')}
              description={t('confirmDeleteDesc')}
              onConfirm={async () => {
                await Promise.all(rows.map((r) => deleteNode({ id: r.id } as any)));
                toast.success(t('deleted'));
                ref.current?.refresh();
                fetchNodes();
                fetchTags();
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
          resetSortWithNode({
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
  );
}
