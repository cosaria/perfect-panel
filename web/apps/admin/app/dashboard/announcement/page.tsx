'use client';

import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createAnnouncement,
  deleteAnnouncement,
  getAnnouncementList,
  updateAnnouncement,
} from '@/services/admin/announcement';
import { Button } from '@workspace/ui/components/button';
import { Switch } from '@workspace/ui/components/switch';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { format } from 'date-fns';
import { useTranslations } from 'next-intl';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import NoticeForm from './notice-form';

export default function Page() {
  const t = useTranslations('announcement');
  const [loading, setLoading] = useState(false);
  const ref = useRef<ProTableActions>(null);

  return (
    <ProTable<API.Announcement, { enable: boolean; search: string }>
      action={ref}
      header={{
        title: t('announcementList'),
        toolbar: (
          <NoticeForm<API.CreateAnnouncementRequest>
            trigger={t('create')}
            title={t('createAnnouncement')}
            loading={loading}
            onSubmit={async (values) => {
              setLoading(true);
              try {
                await createAnnouncement(values);
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
          accessorKey: 'show',
          header: t('show'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('show')}
                onCheckedChange={async (checked) => {
                  await updateAnnouncement({
                    ...row.original,
                    show: checked,
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'pinned',
          header: t('pinned'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('pinned')}
                onCheckedChange={async (checked) => {
                  await updateAnnouncement({
                    ...row.original,
                    pinned: checked,
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'popup',
          header: t('popup'),
          cell: ({ row }) => {
            return (
              <Switch
                defaultChecked={row.getValue('popup')}
                onCheckedChange={async (checked) => {
                  await updateAnnouncement({
                    ...row.original,
                    popup: checked,
                  });
                  ref.current?.refresh();
                }}
              />
            );
          },
        },
        {
          accessorKey: 'title',
          header: t('title'),
        },
        {
          accessorKey: 'content',
          header: t('content'),
        },
        {
          accessorKey: 'updated_at',
          header: t('updatedAt'),
          cell: ({ row }) => format(row.getValue('updated_at'), 'yyyy-MM-dd HH:mm:ss'),
        },
      ]}
      params={[
        {
          key: 'enable',
          placeholder: t('enable'),
          options: [
            { label: t('show'), value: 'false' },
            { label: t('hide'), value: 'true' },
          ],
        },
        { key: 'search' },
      ]}
      request={async (pagination, filter) => {
        const { data } = await getAnnouncementList({
          ...pagination,
          ...filter,
        });
        return {
          list: data.data?.list || [],
          total: data.data?.total || 0,
        };
      }}
      actions={{
        render(row) {
          return [
            <NoticeForm<API.Announcement>
              key='edit'
              trigger={t('edit')}
              title={t('editAnnouncement')}
              loading={loading}
              initialValues={row}
              onSubmit={async (values) => {
                setLoading(true);
                try {
                  await updateAnnouncement({
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
              description={t('deleteDescription')}
              onConfirm={async () => {
                await deleteAnnouncement({
                  id: row.id,
                });
                toast.success(t('deleteSuccess'));
                ref.current?.refresh();
              }}
              cancelText={t('cancel')}
              confirmText={t('confirm')}
            />,
          ];
        },
        batchRender(rows) {
          return [
            <ConfirmButton
              key='delete'
              trigger={<Button variant='destructive'>{t('delete')}</Button>}
              title={t('confirmDelete')}
              description={t('deleteDescription')}
              onConfirm={async () => {
                for (const element of rows) {
                  await deleteAnnouncement({
                    id: element.id!,
                  });
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
