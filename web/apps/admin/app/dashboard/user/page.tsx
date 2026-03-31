'use client';

import { Display } from '@/components/display';
import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createUser,
  deleteUser,
  getUserDetail,
  getUserList,
  updateUserBasicInfo,
} from '@/services/admin/user';
import { useSubscribe } from '@/store/subscribe';
import { formatDate } from '@/utils/common';
import { useQuery } from '@tanstack/react-query';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@workspace/ui/components/dropdown-menu';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { Switch } from '@workspace/ui/components/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import { UserDetail } from './user-detail';
import UserForm from './user-form';
import { AuthMethodsForm } from './user-profile/auth-methods-form';
import { BasicInfoForm } from './user-profile/basic-info-form';
import { NotifySettingsForm } from './user-profile/notify-settings-form';
import UserSubscription from './user-subscription';

export default function Page() {
  const t = useTranslations('user');
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);
  const sp = useSearchParams();

  const { subscribes } = useSubscribe();

  const initialFilters = {
    search: sp.get('search') || undefined,
    user_id: sp.get('user_id') || undefined,
    subscribe_id: sp.get('subscribe_id') || undefined,
    user_subscribe_id: sp.get('user_subscribe_id') || undefined,
  };

  return (
    <ProTable<API.User, API.GetUserListParams>
      key={initialFilters.user_id}
      action={ref}
      initialFilters={initialFilters}
      header={{
        title: t('userList'),
        toolbar: (
          <UserForm<API.CreateUserRequest>
            key='create'
            trigger={t('create')}
            title={t('createUser')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createUser(values);
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
                defaultChecked={row.getValue('enable')}
                onCheckedChange={async (checked) => {
                  const {
                    auth_methods,
                    user_devices,
                    enable_balance_notify,
                    enable_login_notify,
                    enable_subscribe_notify,
                    enable_trade_notify,
                    updated_at,
                    created_at,
                    id,
                    ...rest
                  } = row.original;
                  await updateUserBasicInfo({
                    user_id: id,
                    ...rest,
                    enable: checked,
                  } as unknown as API.UpdateUserBasiceInfoRequest);
                  toast.success(t('updateSuccess'));
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'id',
          header: 'ID',
        },
        {
          accessorKey: 'auth_methods',
          header: t('userName'),
          cell: ({ row }) => {
            const method = row.original.auth_methods?.[0];
            return (
              <div>
                <Badge className='mr-1 uppercase' title={method?.verified ? t('verified') : ''}>
                  {method?.auth_type}
                </Badge>
                {method?.auth_identifier}
              </div>
            );
          },
        },
        {
          accessorKey: 'balance',
          header: t('balance'),
          cell: ({ row }) => <Display type='currency' value={row.getValue('balance')} />,
        },
        {
          accessorKey: 'gift_amount',
          header: t('giftAmount'),
          cell: ({ row }) => <Display type='currency' value={row.getValue('gift_amount')} />,
        },
        {
          accessorKey: 'commission',
          header: t('commission'),
          cell: ({ row }) => <Display type='currency' value={row.getValue('commission')} />,
        },
        {
          accessorKey: 'refer_code',
          header: t('inviteCode'),
          cell: ({ row }) => row.getValue('refer_code') || '--',
        },
        {
          accessorKey: 'referer_id',
          header: t('referer'),
          cell: ({ row }) => <UserDetail id={row.original.referer_id} />,
        },
        {
          accessorKey: 'created_at',
          header: t('createdAt'),
          cell: ({ row }) => formatDate(row.getValue('created_at')),
        },
      ]}
      request={async (pagination, filter) => {
        const { data } = await getUserList({
          ...pagination,
          ...filter,
        });
        return {
          list: data.data?.list || [],
          total: data.data?.total || 0,
        };
      }}
      params={[
        {
          key: 'subscribe_id',
          placeholder: t('subscription'),
          options: subscribes?.map((item) => ({
            label: item.name!,
            value: String(item.id!),
          })),
        },
        {
          key: 'search',
          placeholder: 'Search',
        },
        {
          key: 'user_id',
          placeholder: t('userId'),
        },
        {
          key: 'user_subscribe_id',
          placeholder: t('subscriptionId'),
        },
      ]}
      actions={{
        render: (row) => {
          return [
            <ProfileSheet key='profile' userId={row.id} />,
            <SubscriptionSheet key='subscription' userId={row.id} />,
            <ConfirmButton
              key='edit'
              trigger={<Button variant='destructive'>{t('delete')}</Button>}
              title={t('confirmDelete')}
              description={t('deleteDescription')}
              onConfirm={async () => {
                await deleteUser({ id: row.id });
                toast.success(t('deleteSuccess'));
                ref.current?.refresh();
              }}
              cancelText={t('cancel')}
              confirmText={t('confirm')}
            />,
            <DropdownMenu key='more'>
              <DropdownMenuTrigger asChild>
                <Button variant='outline'>{t('more')}</Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align='end'>
                <DropdownMenuItem asChild>
                  <Link href={`/dashboard/order?user_id=${row.id}`}>{t('orderList')}</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href={`/dashboard/log/login?user_id=${row.id}`}>{t('loginLogs')}</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href={`/dashboard/log/balance?user_id=${row.id}`}>{t('balanceLogs')}</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href={`/dashboard/log/commission?user_id=${row.id}`}>
                    {t('commissionLogs')}
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href={`/dashboard/log/gift?user_id=${row.id}`}>{t('giftLogs')}</Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>,
          ];
        },
      }}
    />
  );
}

function ProfileSheet({ userId }: { userId: number }) {
  const t = useTranslations('user');
  const [open, setOpen] = useState(false);
  const { data: user, refetch } = useQuery({
    enabled: open,
    queryKey: ['user', userId],
    queryFn: async () => {
      const { data } = await getUserDetail({ id: userId });
      return data.data as API.User;
    },
  });
  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant='default'>{t('edit')}</Button>
      </SheetTrigger>
      <SheetContent side='right' className='w-[700px] max-w-full md:max-w-screen-lg'>
        <SheetHeader>
          <SheetTitle>
            {t('userProfile')} · ID: {userId}
          </SheetTitle>
        </SheetHeader>
        {user && (
          <ScrollArea className='h-[calc(100dvh-140px)] p-2'>
            <Tabs defaultValue='basic'>
              <TabsList className='mb-3'>
                <TabsTrigger value='basic'>{t('basicInfoTitle')}</TabsTrigger>
                <TabsTrigger value='notify'>{t('notifySettingsTitle')}</TabsTrigger>
                <TabsTrigger value='auth'>{t('authMethodsTitle')}</TabsTrigger>
              </TabsList>
              <TabsContent value='basic' className='mt-0'>
                <BasicInfoForm user={user} refetch={refetch as any} />
              </TabsContent>
              <TabsContent value='notify' className='mt-0'>
                <NotifySettingsForm user={user} refetch={refetch as any} />
              </TabsContent>
              <TabsContent value='auth' className='mt-0'>
                <AuthMethodsForm user={user} refetch={refetch as any} />
              </TabsContent>
            </Tabs>
          </ScrollArea>
        )}
      </SheetContent>
    </Sheet>
  );
}

function SubscriptionSheet({ userId }: { userId: number }) {
  const t = useTranslations('user');
  const [open, setOpen] = useState(false);
  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant='secondary'>订阅</Button>
      </SheetTrigger>
      <SheetContent side='right' className='w-[1000px] max-w-full md:max-w-screen-xl'>
        <SheetHeader>
          <SheetTitle>
            {t('subscriptionList')} · ID: {userId}
          </SheetTitle>
        </SheetHeader>
        <div className='mt-2'>
          <UserSubscription userId={userId} />
        </div>
      </SheetContent>
    </Sheet>
  );
}
