'use client';

import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createTicketFollow,
  getTicket,
  getTicketList,
  updateTicketStatus,
} from '@/services/admin/ticket';
import { formatDate } from '@/utils/common';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import {
  Drawer,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from '@workspace/ui/components/drawer';
import { Input } from '@workspace/ui/components/input';
import { Label } from '@workspace/ui/components/label';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { Icon } from '@workspace/ui/custom-components/icon';
import { cn } from '@workspace/ui/lib/utils';
import { useTranslations } from 'next-intl';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { UserDetail } from '../user/user-detail';

export default function Page() {
  const t = useTranslations('ticket');

  const [ticketId, setTicketId] = useState<any>(null);

  const [message, setMessage] = useState('');

  const { data: ticket, refetch: refetchTicket } = useQuery({
    queryKey: ['getTicket', ticketId],
    queryFn: async () => {
      const { data } = await getTicket({
        id: ticketId,
      });
      return data.data as API.Ticket;
    },
    enabled: !!ticketId,
    refetchInterval: 5000,
  });

  const scrollRef = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    setTimeout(() => {
      if (scrollRef.current) {
        scrollRef.current.children[1]?.scrollTo({
          top: scrollRef.current.children[1].scrollHeight,
          behavior: 'smooth',
        });
      }
    }, 66);
  }, [ticket?.follow?.length]);

  const ref = useRef<ProTableActions>(null);
  return (
    <>
      <ProTable<API.Ticket, { status: number }>
        action={ref}
        header={{
          title: t('ticketList'),
        }}
        columns={[
          {
            accessorKey: 'title',
            header: t('title'),
          },
          {
            accessorKey: 'user_id',
            header: t('user'),
            cell: ({ row }) => <UserDetail id={row.original.user_id} />,
          },
          {
            accessorKey: 'status',
            header: t('status.0'),
            cell: ({ row }) => (
              <span
                className={cn(
                  'flex items-center gap-2 before:block before:size-1.5 before:animate-pulse before:rounded-full before:ring-2 before:ring-opacity-50',
                  {
                    'before:bg-rose-500 before:ring-rose-500': row.original.status === 1,
                    'before:bg-yellow-500 before:ring-yellow-500': row.original.status === 2,
                    'before:bg-green-500 before:ring-green-500': row.original.status === 3,
                    'before:bg-zinc-500 before:ring-zinc-500': row.original.status === 4,
                  },
                )}
              >
                {t(`status.${row.original.status}`)}
              </span>
            ),
          },
          {
            accessorKey: 'updated_at',
            header: t('updatedAt'),
            cell: ({ row }) => formatDate(row.getValue('updated_at')),
          },
        ]}
        params={[
          {
            key: 'status',
            placeholder: t('status.0'),
            options: [
              {
                label: t('close'),
                value: '4',
              },
            ],
          },
        ]}
        request={async (pagination, filters) => {
          const { data } = await getTicketList({
            ...pagination,
            ...filters,
          });
          return {
            list: data.data?.list || [],
            total: data.data?.total || 0,
          };
        }}
        actions={{
          render(row) {
            if (row.status !== 4) {
              return [
                <Button key='reply' onClick={() => setTicketId(row.id)}>
                  {t('reply')}
                </Button>,
                <ConfirmButton
                  key='colse'
                  trigger={<Button variant='destructive'>{t('close')}</Button>}
                  title={t('confirmClose')}
                  description={t('closeWarning')}
                  onConfirm={async () => {
                    await updateTicketStatus({
                      id: row.id,
                      status: 4,
                    });
                    toast.success(t('closeSuccess'));
                    ref.current?.refresh();
                  }}
                  cancelText={t('cancel')}
                  confirmText={t('confirm')}
                />,
              ];
            }
            return [
              <Button key='check' size='sm' onClick={() => setTicketId(row.id)}>
                {t('check')}
              </Button>,
            ];
          },
        }}
      />

      <Drawer
        open={!!ticketId}
        onOpenChange={(open) => {
          if (!open) setTicketId(null);
        }}
      >
        <DrawerContent className='container mx-auto h-screen *:select-text'>
          <DrawerHeader className='border-b text-left'>
            <DrawerTitle>{ticket?.title}</DrawerTitle>
          </DrawerHeader>
          <ScrollArea className='h-full overflow-hidden' ref={scrollRef}>
            <div className='flex h-full flex-col gap-4 p-4'>
              {/* 显示工单描述作为第一条用户消息 */}
              {ticket?.description && (
                <div className='flex items-center gap-4'>
                  <div className='flex flex-col gap-1'>
                    <p className='text-muted-foreground text-sm'>{formatDate(ticket.created_at)}</p>
                    <p className='bg-accent w-fit rounded-lg p-2 font-medium'>
                      {ticket.description}
                    </p>
                  </div>
                </div>
              )}

              {/* 显示后续跟进消息 */}
              {ticket?.follow?.map((item) => (
                <div
                  key={item.id}
                  className={cn('flex items-center gap-4', {
                    'flex-row-reverse': item.from === 'System',
                  })}
                >
                  <div
                    className={cn('flex flex-col gap-1', {
                      'items-end': item.from === 'System',
                    })}
                  >
                    <p className='text-muted-foreground text-sm'>{formatDate(item.created_at)}</p>
                    <p
                      className={cn('bg-accent w-fit rounded-lg p-2 font-medium', {
                        'bg-primary text-primary-foreground': item.from === 'System',
                      })}
                    >
                      {item.type === 1 && item.content}
                      {item.type === 2 && (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          src={item.content!}
                          width={300}
                          height={300}
                          className='!size-auto object-cover'
                          alt='image'
                        />
                      )}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>
          {ticket?.status !== 4 && (
            <DrawerFooter>
              <form
                className='flex w-full flex-row items-center gap-2'
                onSubmit={async (event) => {
                  event.preventDefault();
                  if (message) {
                    await createTicketFollow({
                      ticket_id: ticketId,
                      from: 'System',
                      type: 1,
                      content: message,
                    });
                    refetchTicket();
                    setMessage('');
                  }
                }}
              >
                <Button type='button' variant='outline' className='p-0'>
                  <Label htmlFor='picture' className='p-2'>
                    <Icon icon='uil:image-upload' className='text-2xl' />
                  </Label>
                  <Input
                    id='picture'
                    type='file'
                    className='hidden'
                    accept='image/*'
                    onChange={(event) => {
                      const file = event.target.files?.[0];
                      if (file && file.type.startsWith('image/')) {
                        const reader = new FileReader();
                        reader.readAsDataURL(file);
                        reader.onload = (e) => {
                          const img = new Image();
                          img.src = e.target?.result as string;
                          img.onload = () => {
                            const canvas = document.createElement('canvas');
                            const ctx = canvas.getContext('2d');

                            const maxWidth = 300;
                            const maxHeight = 300;
                            let width = img.width;
                            let height = img.height;

                            if (width > height) {
                              if (width > maxWidth) {
                                height = Math.round((maxWidth / width) * height);
                                width = maxWidth;
                              }
                            } else {
                              if (height > maxHeight) {
                                width = Math.round((maxHeight / height) * width);
                                height = maxHeight;
                              }
                            }

                            canvas.width = width;
                            canvas.height = height;
                            ctx?.drawImage(img, 0, 0, width, height);

                            canvas.toBlob(
                              (blob) => {
                                const reader = new FileReader();
                                reader.readAsDataURL(blob!);
                                reader.onloadend = async () => {
                                  await createTicketFollow({
                                    ticket_id: ticketId,
                                    from: 'System',
                                    type: 2,
                                    content: reader.result as string,
                                  });
                                  refetchTicket();
                                };
                              },
                              'image/webp',
                              0.8,
                            );
                          };
                        };
                      }
                    }}
                  />
                </Button>
                <Input
                  placeholder={t('inputPlaceholder')}
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                />
                <Button type='submit' disabled={!message}>
                  <Icon icon='uil:navigator' />
                </Button>
              </form>
            </DrawerFooter>
          )}
        </DrawerContent>
      </Drawer>
    </>
  );
}
