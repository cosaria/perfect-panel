'use client';

import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';
import { useRef } from 'react';

import { Display } from '@/components/display';
import { ProTable, ProTableActions } from '@/components/pro-table';
import { getOrderList, updateOrderStatus } from '@/services/admin/order';
import { useSubscribe } from '@/store/subscribe';
import { formatDate } from '@/utils/common';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@workspace/ui/components/hover-card';
import { Separator } from '@workspace/ui/components/separator';
import { Combobox } from '@workspace/ui/custom-components/combobox';
import { cn } from '@workspace/ui/lib/utils';
import { UserDetail } from '../user/user-detail';

export default function Page() {
  const t = useTranslations('order');
  const sp = useSearchParams();

  const statusOptions = [
    { value: 1, label: t('status.1'), className: 'bg-orange-500' },
    { value: 2, label: t('status.2'), className: 'bg-green-500' },
    { value: 3, label: t('status.3'), className: 'bg-gray-500' },
    { value: 4, label: t('status.4'), className: 'bg-red-500' },
    { value: 5, label: t('status.5'), className: 'bg-green-500' },
  ];

  const ref = useRef<ProTableActions>(null);

  const { subscribes, getSubscribeName } = useSubscribe();

  const initialFilters = {
    search: sp.get('search') || undefined,
    status: sp.get('status') || undefined,
    subscribe_id: sp.get('subscribe_id') || undefined,
    user_id: sp.get('user_id') || undefined,
  };

  return (
    <ProTable<API.Order, any>
      action={ref}
      initialFilters={initialFilters}
      columns={[
        {
          accessorKey: 'order_no',
          header: t('orderNumber'),
        },
        {
          accessorKey: 'type',
          header: t('type.0'),
          cell: ({ row }) => t(`type.${row.getValue('type')}`),
        },
        {
          accessorKey: 'subscribe_id',
          header: t('subscribe'),
          cell: ({ row }) => {
            if (row.original.type === 4) {
              return t(`type.${row.getValue('type')}`);
            }
            const name = getSubscribeName(row.getValue('subscribe_id'));
            const quantity = row.original.quantity;
            return name ? `${name} × ${quantity}` : '';
          },
        },
        {
          accessorKey: 'amount',
          header: t('amount'),
          cell: ({ row }) => (
            <HoverCard>
              <HoverCardTrigger asChild>
                <Button variant='link' className='p-0'>
                  <Display type='currency' value={row.original.amount} />
                </Button>
              </HoverCardTrigger>
              <HoverCardContent>
                <div className='grid gap-3'>
                  {row.original.trade_no && (
                    <>
                      <div className='font-semibold'>{t('tradeNo')}</div>
                      <span className='text-muted-foreground'>{row.original.trade_no}</span>
                      <Separator className='my-2' />
                    </>
                  )}
                  <ul className='grid gap-3'>
                    <li className='flex items-center justify-between'>
                      <span className='text-muted-foreground'>{t('subscribePrice')}</span>
                      <span>
                        <Display type='currency' value={row.original.price} />
                      </span>
                    </li>
                    <li className='flex items-center justify-between'>
                      <span className='text-muted-foreground'>{t('discount')}</span>
                      <span>
                        <Display type='currency' value={row.original.discount} />
                      </span>
                    </li>
                    <li className='flex items-center justify-between'>
                      <span className='text-muted-foreground'>{t('couponDiscount')}</span>
                      <span>
                        <Display type='currency' value={row.original.coupon_discount} />
                      </span>
                    </li>
                    <li className='flex items-center justify-between'>
                      <span className='text-muted-foreground'>{t('feeAmount')}</span>
                      <span>
                        <Display type='currency' value={row.original.fee_amount} />
                      </span>
                    </li>
                    <li className='flex items-center justify-between font-semibold'>
                      <span className='text-muted-foreground'>{t('total')}</span>
                      <span>
                        <Display type='currency' value={row.original.amount} />
                      </span>
                    </li>
                  </ul>
                </div>
                <Separator className='my-4' />
                <ul className='grid gap-3'>
                  <li className='flex items-center justify-between'>
                    <span className='text-muted-foreground'>{t('method')}</span>
                    <span>{row.original?.payment?.name || row.original?.payment?.platform}</span>
                  </li>
                </ul>
              </HoverCardContent>
            </HoverCard>
          ),
        },
        {
          accessorKey: 'user_id',
          header: t('user'),
          cell: ({ row }) => <UserDetail id={row.original.user_id} />,
        },
        {
          accessorKey: 'updated_at',
          header: t('updateTime'),
          cell: ({ row }) => formatDate(row.original.updated_at),
        },
        {
          accessorKey: 'status',
          header: t('status.0'),
          cell: ({ row }) => {
            const option = statusOptions.find((opt) => opt.value === row.original.status);
            if ([1, 3, 4].includes(row.getValue('status'))) {
              return (
                <Combobox<number, false>
                  placeholder={t('status.0')}
                  value={row.original.status}
                  onChange={async (value) => {
                    await updateOrderStatus({
                      id: row.original.id,
                      status: value,
                    });
                    ref.current?.refresh();
                  }}
                  options={statusOptions}
                  className={cn(option?.className)}
                />
              );
            }
            return <Badge>{t(`status.${row.getValue('status')}`)}</Badge>;
          },
        },
      ]}
      params={[
        {
          key: 'status',
          placeholder: t('status.0'),
          options: statusOptions.map((item) => ({
            label: item.label,
            value: String(item.value),
          })),
        },
        {
          key: 'subscribe_id',
          placeholder: `${t('subscribe')}`,
          options: subscribes?.map((item) => ({
            label: item.name!,
            value: String(item.id),
          })),
        },
        { key: 'search' },
        {
          key: 'user_id',
          placeholder: `${t('user')} ID`,
          options: undefined,
        },
      ]}
      request={async (pagination, filter) => {
        const { data } = await getOrderList({ ...pagination, ...filter });
        return {
          list: data.data?.list || [],
          total: data.data?.total || 0,
        };
      }}
    />
  );
}
