'use client';

import { Display } from '@/components/display';
import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createPaymentMethod,
  deletePaymentMethod,
  getPaymentMethodList,
  updatePaymentMethod,
} from '@/services/admin/payment';
import { Avatar, AvatarFallback, AvatarImage } from '@workspace/ui/components/avatar';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import { Switch } from '@workspace/ui/components/switch';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import PaymentForm from './payment-form';

export default function PaymentTable() {
  const t = useTranslations('payment');
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);

  return (
    <ProTable<API.PaymentConfig, { search: string }>
      action={ref}
      header={{
        title: t('paymentManagement'),
        toolbar: (
          <PaymentForm<API.CreatePaymentMethodRequest>
            trigger={<Button>{t('create')}</Button>}
            title={t('createPayment')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createPaymentMethod({
                  ...values,
                  enable: false,
                });
                toast.success(t('createSuccess'));
                ref.current?.refresh();
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
      columns={[
        {
          accessorKey: 'enable',
          header: t('enable'),
          cell: ({ row }) => {
            return (
              <Switch
                checked={Boolean(row.getValue('enable'))}
                onCheckedChange={async (checked) => {
                  await updatePaymentMethod({
                    ...row.original,
                    enable: checked,
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'icon',
          header: t('icon'),
          cell: ({ row }) => {
            const icon = row.getValue('icon') as string;
            return (
              <Avatar className='h-8 w-8'>
                {icon ? <AvatarImage src={icon} alt={row.getValue('name')} /> : null}
                <AvatarFallback>
                  {(row.getValue('name') as string)?.charAt(0) || '?'}
                </AvatarFallback>
              </Avatar>
            );
          },
        },
        {
          accessorKey: 'name',
          header: t('name'),
        },
        {
          accessorKey: 'platform',
          header: t('platform'),
          cell: ({ row }) => <Badge>{t(row.original.platform)}</Badge>,
        },
        {
          accessorKey: 'notify_url',
          header: t('notify_url'),
        },
        {
          accessorKey: 'fee',
          header: t('handlingFee'),
          cell: ({ row }) => {
            const feeMode = row.original.fee_mode;
            if (feeMode === 1) {
              return <Badge>{row.original.fee_percent}%</Badge>;
            } else if (feeMode === 2) {
              return (
                <Badge>
                  <Display value={row.original.fee_amount} type='currency' />
                </Badge>
              );
            }
            return '--';
          },
        },
      ]}
      params={[
        {
          key: 'search',
          placeholder: t('searchPlaceholder'),
        },
      ]}
      request={async (pagination, filter) => {
        const { data } = await getPaymentMethodList({
          ...pagination,
          ...filter,
        });
        return {
          list: data?.data?.list || [],
          total: data?.data?.total || 0,
        };
      }}
      actions={{
        render: (row) => [
          <PaymentForm<API.UpdatePaymentMethodRequest>
            isEdit
            key='edit'
            trigger={<Button>{t('edit')}</Button>}
            title={t('editPayment')}
            loading={loading}
            initialValues={row}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await updatePaymentMethod({
                  ...row,
                  ...values,
                });
                toast.success(t('updateSuccess'));
                ref.current?.refresh();
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
              await deletePaymentMethod({
                id: row.id,
              });
              toast.success(t('deleteSuccess'));
              ref.current?.refresh();
            }}
            cancelText={t('cancel')}
            confirmText={t('confirm')}
          />,
          <Button
            key='copy'
            variant='outline'
            onClick={async () => {
              setLoading(true);
              try {
                const { id, ...params } = row;
                await createPaymentMethod({
                  ...params,
                  enable: false,
                });
                toast.success(t('copySuccess'));
                ref.current?.refresh();
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
        batchRender(rows) {
          return [
            <ConfirmButton
              key='delete'
              trigger={<Button variant='destructive'>{t('batchDelete')}</Button>}
              title={t('confirmDelete')}
              description={t('deleteWarning')}
              onConfirm={async () => {
                for (const row of rows) {
                  await deletePaymentMethod({ id: row.id });
                }
                toast.success(t('deleteSuccess'));
                ref.current?.refresh();
              }}
              cancelText={t('cancel')}
              confirmText={t('confirm')}
            />,
          ];
        },
      }}
    />
  );
}
