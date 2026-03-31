'use client';

import { Display } from '@/components/display';
import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  batchDeleteSubscribe,
  createSubscribe,
  deleteSubscribe,
  getSubscribeList,
  subscribeSort,
  updateSubscribe,
} from '@/services/admin/subscribe';
import { useSubscribe } from '@/store/subscribe';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Switch } from '@workspace/ui/components/switch';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import SubscribeForm from './subscribe-form';

export default function SubscribeTable() {
  const t = useTranslations('product');
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);
  const { fetchSubscribes } = useSubscribe();
  return (
    <ProTable<API.SubscribeItem, { group_id: number; query: string }>
      action={ref}
      header={{
        toolbar: (
          <SubscribeForm<API.CreateSubscribeRequest>
            trigger={t('create')}
            title={t('createSubscribe')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createSubscribe({
                  ...values,
                  show: false,
                  sell: false,
                });
                toast.success(t('createSuccess'));
                ref.current?.refresh();
                fetchSubscribes();
                setLoading(false);

                return true;
              } catch (error) {
                setLoading(false);

                return false;
              }
            }}
          />
        ),
      }}
      params={[
        {
          key: 'search',
        },
      ]}
      request={async (pagination, filters) => {
        const { data } = await getSubscribeList({
          ...pagination,
          ...filters,
        });
        return {
          list: data.data?.list || [],
          total: data.data?.total || 0,
        };
      }}
      columns={[
        {
          accessorKey: 'show',
          header: t('show'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('show')}
                onCheckedChange={async (checked) => {
                  await updateSubscribe({
                    ...row.original,
                    show: checked,
                  } as API.UpdateSubscribeRequest);
                  ref.current?.refresh();
                  fetchSubscribes();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'sell',
          header: t('sell'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('sell')}
                onCheckedChange={async (checked) => {
                  await updateSubscribe({
                    ...row.original,
                    sell: checked,
                  } as API.UpdateSubscribeRequest);
                  ref.current?.refresh();
                  fetchSubscribes();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'name',
          header: t('name'),
        },
        {
          accessorKey: 'unit_price',
          header: t('unitPrice'),
          cell: ({ row }) => {
            return (
              <>
                <Display type='currency' value={row.getValue('unit_price')} />/
                {t(row.original.unit_time ? `form.${row.original.unit_time}` : 'form.Month')}
              </>
            );
          },
        },
        {
          accessorKey: 'replacement',
          header: t('replacement'),
          cell: ({ row }) => <Display type='currency' value={row.getValue('replacement')} />,
        },
        {
          accessorKey: 'traffic',
          header: t('traffic'),
          cell: ({ row }) => <Display type='traffic' value={row.getValue('traffic')} unlimited />,
        },
        {
          accessorKey: 'device_limit',
          header: t('deviceLimit'),
          cell: ({ row }) => (
            <Display type='number' value={row.getValue('device_limit')} unlimited />
          ),
        },
        {
          accessorKey: 'inventory',
          header: t('inventory'),
          cell: ({ row }) => (
            <Display
              type='number'
              value={row.getValue('inventory') === -1 ? 0 : row.getValue('inventory')}
              unlimited
            />
          ),
        },
        {
          accessorKey: 'quota',
          header: t('quota'),
          cell: ({ row }) => <Display type='number' value={row.getValue('quota')} unlimited />,
        },
        {
          accessorKey: 'language',
          header: t('language'),
          cell: ({ row }) => {
            const language = row.getValue('language') as string;
            return language ? <Badge variant='outline'>{language}</Badge> : '--';
          },
        },
        {
          accessorKey: 'sold',
          header: t('sold'),
          cell: ({ row }) => <Badge variant='outline'>{row.getValue('sold')}</Badge>,
        },
      ]}
      actions={{
        render: (row) => [
          <SubscribeForm<API.SubscribeItem>
            key='edit'
            trigger={t('edit')}
            title={t('editSubscribe')}
            loading={loading}
            initialValues={row}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await updateSubscribe({
                  ...row,
                  ...values,
                } as API.UpdateSubscribeRequest);
                toast.success(t('updateSuccess'));
                ref.current?.refresh();
                fetchSubscribes();
                setLoading(false);
                return true;
              } catch (error) {
                setLoading(false);

                return false;
              }
            }}
          />,
          <ConfirmButton
            key='delete'
            trigger={<Button variant='destructive'>{t('delete')}</Button>}
            title={t('confirmDelete')}
            description={t('deleteWarning')}
            onConfirm={async () => {
              await deleteSubscribe({
                id: row.id!,
              });
              toast.success(t('deleteSuccess'));
              ref.current?.refresh();
              fetchSubscribes();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
          <Button
            key='copy'
            variant='secondary'
            onClick={async () => {
              setLoading(true);
              try {
                const { id, sort, sell, updated_at, created_at, ...params } = row;
                await createSubscribe({
                  ...params,
                  show: false,
                  sell: false,
                } as API.CreateSubscribeRequest);
                toast.success(t('copySuccess'));
                ref.current?.refresh();
                fetchSubscribes();
                setLoading(false);
                return true;
              } catch (error) {
                setLoading(false);
                return false;
              }
            }}
          >
            {t('copy')}
          </Button>,
        ],
        batchRender: (rows) => [
          <ConfirmButton
            key='delete'
            trigger={<Button variant='destructive'>{t('delete')}</Button>}
            title={t('confirmDelete')}
            description={t('deleteWarning')}
            onConfirm={async () => {
              await batchDeleteSubscribe({
                ids: rows.map((item) => item.id) as number[],
              });

              toast.success(t('deleteSuccess'));
              ref.current?.reset();
              fetchSubscribes();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
        ],
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
          subscribeSort({
            sort: changedItems.map((item) => ({ id: item.id, sort: item.sort })) as API.SortItem[],
          });
        }

        return updatedItems;
      }}
    />
  );
}
