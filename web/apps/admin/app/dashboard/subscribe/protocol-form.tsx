'use client';

import { ProTable, ProTableActions } from '@/components/pro-table';
import {
  createSubscribeApplication,
  deleteSubscribeApplication,
  getSubscribeApplicationList,
  updateSubscribeApplication,
} from '@/services/admin/application';
import { zodResolver } from '@hookform/resolvers/zod';
import { ColumnDef } from '@tanstack/react-table';
import { Badge } from '@workspace/ui/components/badge';
import { Button } from '@workspace/ui/components/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@workspace/ui/components/form';
import { Input } from '@workspace/ui/components/input';
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@workspace/ui/components/select';
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@workspace/ui/components/sheet';
import { Switch } from '@workspace/ui/components/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import { Textarea } from '@workspace/ui/components/textarea';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@workspace/ui/components/tooltip';
import { ConfirmButton } from '@workspace/ui/custom-components/confirm-button';
import { GoTemplateEditor } from '@workspace/ui/custom-components/editor';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { UploadImage } from '@workspace/ui/custom-components/upload-image';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import { useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
import { subscribeSchema } from './schema';
import { TemplatePreview } from './template-preview';

const createClientFormSchema = (t: any) =>
  z.object({
    name: z.string().min(1, t('form.validation.nameRequired')),
    description: z.string().optional(),
    icon: z.string().optional(),
    user_agent: z.string().min(1, `User-Agent ${t('form.validation.userAgentRequiredSuffix')}`),
    scheme: z.string().optional(),
    template: z.string(),
    output_format: z.string(),
    download_link: z.object({
      windows: z.string().optional(),
      mac: z.string().optional(),
      linux: z.string().optional(),
      ios: z.string().optional(),
      android: z.string().optional(),
      harmony: z.string().optional(),
    }),
  });

type ClientFormData = z.infer<ReturnType<typeof createClientFormSchema>>;

export function ProtocolForm() {
  const t = useTranslations('subscribe');
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [editingClient, setEditingClient] = useState<API.SubscribeApplication | null>(null);
  const tableRef = useRef<ProTableActions>(null);

  const clientFormSchema = createClientFormSchema(t);

  const form = useForm<ClientFormData>({
    resolver: zodResolver(clientFormSchema),
    defaultValues: {
      name: '',
      description: '',
      icon: '',
      user_agent: '',
      scheme: '',
      template: '',
      output_format: '',
      download_link: {
        windows: '',
        mac: '',
        linux: '',
        ios: '',
        android: '',
        harmony: '',
      },
    },
  });

  const request = async (
    pagination: { page: number; size: number },
    filter: Record<string, unknown>,
  ) => {
    const { data } = await getSubscribeApplicationList({
      page: pagination.page,
      size: pagination.size,
    });

    return {
      list: data.data?.list || [],
      total: data.data?.total || 0,
    };
  };

  const columns: ColumnDef<API.SubscribeApplication, any>[] = [
    {
      accessorKey: 'is_default',
      header: t('table.columns.default'),
      cell: ({ row }) => (
        <Switch
          checked={row.original.is_default || false}
          onCheckedChange={async (checked) => {
            await updateSubscribeApplication({
              ...row.original,
              is_default: checked,
            });
            tableRef.current?.refresh();
          }}
        />
      ),
    },
    {
      accessorKey: 'name',
      header: t('table.columns.name'),
      cell: ({ row }) => (
        <div className='flex items-center gap-2'>
          {row.original.icon && (
            <div className='relative h-6 w-6 flex-shrink-0'>
              <Image
                src={row.original.icon}
                alt={row.original.name}
                width={24}
                height={24}
                className='rounded object-contain'
                onError={() => {
                  console.log(`Failed to load image for ${row.original.name}`);
                }}
              />
            </div>
          )}
          <span className='font-medium'>{row.original.name}</span>
        </div>
      ),
    },
    {
      accessorKey: 'user_agent',
      header: 'User-Agent',
      cell: ({ row }) => (
        <div className='text-muted-foreground max-w-[150px] truncate font-mono text-sm'>
          {row.original.user_agent}
        </div>
      ),
    },
    {
      accessorKey: 'output_format',
      header: t('table.columns.outputFormat'),
      cell: ({ row }) => (
        <Badge variant='secondary' className='text-xs'>
          {t(`outputFormats.${row.original.output_format}`) || row.original.output_format}
        </Badge>
      ),
    },
    {
      accessorKey: 'download_link',
      header: t('table.columns.supportedPlatforms'),
      cell: ({ row }) => {
        return (
          <div className='flex flex-wrap gap-1'>
            {Object.entries(row.original.download_link || {}).map(([key, value]) => {
              if (value) {
                return (
                  <Badge key={key} variant='secondary' className='text-xs'>
                    {t(`platforms.${key}`)}
                  </Badge>
                );
              }
              return null;
            })}
          </div>
        );
      },
    },
    {
      accessorKey: 'description',
      header: t('table.columns.description'),
      cell: ({ row }) => (
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>
              <div className='text-muted-foreground max-w-[200px] truncate text-sm'>
                {row.original.description}
              </div>
            </TooltipTrigger>
            <TooltipContent className='max-w-xs'>
              <p>{row.original.description}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      ),
    },
  ];

  const handleAdd = () => {
    setEditingClient(null);
    form.reset({
      name: '',
      description: '',
      icon: '',
      user_agent: '',
      scheme: '',
      template: '',
      output_format: '',
      download_link: {
        windows: '',
        mac: '',
        linux: '',
        ios: '',
        android: '',
        harmony: '',
      },
    });
    setOpen(true);
  };

  const handleEdit = (client: API.SubscribeApplication) => {
    setEditingClient(client);
    form.reset(client);
    setOpen(true);
  };

  const handleDelete = async (client: API.SubscribeApplication) => {
    setLoading(true);
    try {
      await deleteSubscribeApplication({ id: client.id });
      tableRef.current?.refresh();
      toast.success(t('actions.deleteSuccess'));
    } catch (error) {
      console.error('Failed to delete client:', error);
      toast.error(t('actions.deleteFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleBatchDelete = async (clients: API.SubscribeApplication[]) => {
    setLoading(true);
    try {
      await Promise.all(clients.map((client) => deleteSubscribeApplication({ id: client.id })));
      tableRef.current?.refresh();
      toast.success(t('actions.batchDeleteSuccess', { count: clients.length }));
    } catch (error) {
      console.error('Failed to batch delete clients:', error);
      toast.error(t('actions.deleteFailed'));
    } finally {
      setLoading(false);
    }
  };

  const onSubmit = async (data: ClientFormData) => {
    setLoading(true);
    try {
      if (editingClient) {
        await updateSubscribeApplication({
          ...data,
          is_default: editingClient.is_default,
          id: editingClient.id,
        });
        toast.success(t('actions.updateSuccess'));
      } else {
        await createSubscribeApplication({
          ...data,
          is_default: false,
        });
        toast.success(t('actions.createSuccess'));
      }

      setOpen(false);
      tableRef.current?.refresh();
    } catch (error) {
      console.error('Failed to save client:', error);
      toast.error(t('actions.saveFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <ProTable<API.SubscribeApplication, Record<string, unknown>>
        action={tableRef}
        columns={columns}
        request={request}
        header={{
          title: (
            <div className='flex items-center justify-between'>
              <h2 className='text-lg font-semibold'>{t('protocol.title')}</h2>
              <a
                href='https://github.com/perfect-panel/subscription-template'
                target='_blank'
                rel='noreferrer'
                className='text-primary inline-flex items-center gap-2 rounded-md px-3 py-1 text-sm font-medium hover:underline'
              >
                <Icon icon='mdi:github' className='h-4 w-4' />
                <span>Template Repo</span>
                <Icon icon='mdi:open-in-new' className='text-muted-foreground h-4 w-4' />
              </a>
            </div>
          ),
          toolbar: <Button onClick={handleAdd}>{t('actions.add')}</Button>,
        }}
        actions={{
          render: (row) => [
            <TemplatePreview
              key='preview'
              applicationId={row.id}
              output_format={row.output_format}
            />,
            <Button
              key='edit'
              onClick={() => handleEdit(row as unknown as API.SubscribeApplication)}
            >
              {t('actions.edit')}
            </Button>,
            <ConfirmButton
              key='delete'
              trigger={
                <Button variant='destructive' disabled={loading}>
                  {t('actions.delete')}
                </Button>
              }
              title={t('actions.confirmDelete')}
              description={t('actions.deleteWarning')}
              onConfirm={() => handleDelete(row as unknown as API.SubscribeApplication)}
              cancelText={t('actions.cancel')}
              confirmText={t('actions.confirm')}
            />,
          ],
          batchRender: (rows) => [
            <ConfirmButton
              key='batchDelete'
              trigger={<Button variant='destructive'>{t('actions.batchDelete')}</Button>}
              title={t('actions.confirmDelete')}
              description={t('actions.batchDeleteWarning', { count: rows.length })}
              onConfirm={() => handleBatchDelete(rows as unknown as API.SubscribeApplication[])}
              cancelText={t('actions.cancel')}
              confirmText={t('actions.confirm')}
            />,
          ],
        }}
      />

      <Sheet open={open} onOpenChange={setOpen}>
        <SheetContent className='w-[580px] max-w-full md:max-w-screen-md'>
          <SheetHeader>
            <SheetTitle>{editingClient ? t('form.editTitle') : t('form.addTitle')}</SheetTitle>
          </SheetHeader>
          <ScrollArea className='h-[calc(100dvh-48px-36px-36px)]'>
            <Form {...form}>
              <form className='space-y-6 py-4'>
                <Tabs defaultValue='basic' className='w-full'>
                  <TabsList className='grid w-full grid-cols-3'>
                    <TabsTrigger value='basic'>{t('form.tabs.basic')}</TabsTrigger>
                    <TabsTrigger value='template'>{t('form.tabs.template')}</TabsTrigger>
                    <TabsTrigger value='download'>{t('form.tabs.download')}</TabsTrigger>
                  </TabsList>

                  <TabsContent value='basic' className='space-y-4'>
                    <FormField
                      control={form.control}
                      name='icon'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('form.fields.icon')}</FormLabel>
                          <FormControl>
                            <EnhancedInput
                              type='text'
                              placeholder='https://example.com/icon.png'
                              suffix={
                                <UploadImage
                                  className='bg-muted h-9 rounded-none border-none px-2'
                                  onChange={(value) => {
                                    form.setValue(field.name, value as string);
                                  }}
                                />
                              }
                              value={field.value}
                              onValueChange={(value) => {
                                form.setValue(field.name, value as string);
                              }}
                            />
                          </FormControl>
                          <FormDescription>{t('form.descriptions.icon')}</FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name='name'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('form.fields.name')}</FormLabel>
                          <FormControl>
                            <Input placeholder='Clash for Windows' {...field} />
                          </FormControl>
                          <FormDescription>{t('form.descriptions.name')}</FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name='user_agent'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>User-Agent</FormLabel>
                          <FormControl>
                            <Input placeholder='Clash' {...field} />
                          </FormControl>
                          <FormDescription>
                            {t('form.descriptions.userAgentPrefix')}
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name='description'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('form.fields.description')}</FormLabel>
                          <FormControl>
                            <Textarea placeholder={t('form.descriptions.description')} {...field} />
                          </FormControl>
                          <FormDescription>{t('form.descriptions.description')}</FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </TabsContent>

                  <TabsContent value='template' className='space-y-4'>
                    <FormField
                      control={form.control}
                      name='output_format'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('form.fields.outputFormat')}</FormLabel>
                          <FormControl>
                            <Select value={field.value} onValueChange={field.onChange}>
                              <SelectTrigger>
                                <SelectValue placeholder='Select ...' />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value='base64'>{t('outputFormats.base64')}</SelectItem>
                                <SelectItem value='yaml'>{t('outputFormats.yaml')}</SelectItem>
                                <SelectItem value='json'>{t('outputFormats.json')}</SelectItem>
                                <SelectItem value='conf'>{t('outputFormats.conf')}</SelectItem>
                                <SelectItem value='plain'>{t('outputFormats.plain')}</SelectItem>
                              </SelectContent>
                            </Select>
                          </FormControl>
                          <FormDescription>{t('form.descriptions.outputFormat')}</FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name='scheme'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className='flex items-center gap-2'>
                            {t('form.fields.scheme')}
                            <TooltipProvider>
                              <Tooltip>
                                <TooltipTrigger
                                  onClick={(e) => {
                                    e.preventDefault();
                                    e.stopPropagation();
                                  }}
                                >
                                  <Icon
                                    icon='mdi:help-circle-outline'
                                    className='text-muted-foreground h-4 w-4'
                                  />
                                </TooltipTrigger>
                                <TooltipContent
                                  side='right'
                                  className='bg-secondary text-secondary-foreground max-w-md'
                                >
                                  <div className='space-y-2 text-sm'>
                                    <div className='font-medium'>
                                      {t('form.descriptions.scheme.title')}
                                    </div>

                                    <div>
                                      <div className='font-medium'>
                                        {t('form.descriptions.scheme.variables')}
                                      </div>
                                      <ul className='ml-2 list-disc space-y-1 text-xs'>
                                        <li>
                                          <code className='rounded px-1'>{'${url}'}</code> -{' '}
                                          {t('form.descriptions.scheme.urlVariable')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>{'${name}'}</code> -{' '}
                                          {t('form.descriptions.scheme.nameVariable')}
                                        </li>
                                      </ul>
                                    </div>

                                    <div>
                                      <div className='font-medium'>
                                        {t('form.descriptions.scheme.functions')}
                                      </div>
                                      <ul className='ml-2 list-disc space-y-1 text-xs'>
                                        <li>
                                          <code className='rounded px-1'>
                                            {'${encodeURIComponent(...)}'}
                                          </code>{' '}
                                          - {t('form.descriptions.scheme.urlEncoding')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>
                                            {'${window.btoa(...)}'}
                                          </code>{' '}
                                          - {t('form.descriptions.scheme.base64Encoding')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>
                                            {'${JSON.stringify(...)}'}
                                          </code>{' '}
                                          - {t('form.descriptions.scheme.jsonStringify')}
                                        </li>
                                      </ul>
                                    </div>
                                  </div>
                                </TooltipContent>
                              </Tooltip>
                            </TooltipProvider>
                          </FormLabel>
                          <FormControl>
                            <Input
                              placeholder='clash://install-config?url=${url}&name=${name}'
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name='template'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className='flex items-center gap-2'>
                            {t('form.fields.template')}
                            <TooltipProvider>
                              <Tooltip>
                                <TooltipTrigger
                                  onClick={(e) => {
                                    e.preventDefault();
                                    e.stopPropagation();
                                  }}
                                >
                                  <Icon
                                    icon='mdi:help-circle-outline'
                                    className='text-muted-foreground h-4 w-4'
                                  />
                                </TooltipTrigger>
                                <TooltipContent
                                  side='right'
                                  className='bg-secondary text-secondary-foreground max-w-md'
                                >
                                  <div className='space-y-2 text-sm'>
                                    <div className='font-medium'>
                                      {t('form.descriptions.template.title')}
                                    </div>

                                    <div>
                                      <div className='font-medium'>
                                        {t('form.descriptions.template.variables')}
                                      </div>
                                      <ul className='ml-2 list-disc space-y-1 text-xs'>
                                        <li>
                                          <code className='rounded px-1'>.SiteName</code> -{' '}
                                          {t('form.descriptions.template.siteName')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>.SubscribeName</code> -{' '}
                                          {t('form.descriptions.template.subscribeName')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>.Proxies</code> -{' '}
                                          {t('form.descriptions.template.nodes')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>.UserInfo</code> -{' '}
                                          {t('form.descriptions.template.userInfo')}
                                        </li>
                                      </ul>
                                    </div>

                                    <div>
                                      <div className='font-medium'>
                                        {t('form.descriptions.template.functions')}
                                      </div>
                                      <ul className='ml-2 list-disc space-y-1 text-xs'>
                                        <li>
                                          <code className='rounded px-1'>
                                            {'{{range .Proxies}}...{{end}}'}
                                          </code>{' '}
                                          - {t('form.descriptions.template.range')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>
                                            {'{{if .condition}}...{{end}}'}
                                          </code>{' '}
                                          - {t('form.descriptions.template.if')}
                                        </li>
                                        <li>
                                          <code className='rounded px-1'>{'{{sprig_func}}'}</code> -{' '}
                                          {t('form.descriptions.template.sprig')}
                                        </li>
                                      </ul>
                                    </div>
                                  </div>
                                </TooltipContent>
                              </Tooltip>
                            </TooltipProvider>
                          </FormLabel>
                          <FormControl>
                            <GoTemplateEditor
                              showLineNumbers
                              schema={subscribeSchema}
                              enableSprig
                              value={field.value || ''}
                              onChange={(value) => field.onChange(value)}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </TabsContent>

                  <TabsContent value='download' className='space-y-4'>
                    <div className='space-y-4'>
                      <div className='grid gap-3'>
                        {['windows', 'mac', 'linux', 'ios', 'android', 'harmony'].map((key) => (
                          <FormField
                            key={key}
                            control={form.control}
                            name={
                              `download_link.${key}` as `download_link.${keyof ClientFormData['download_link']}`
                            }
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>{t(`platforms.${key}`)}</FormLabel>
                                <FormControl>
                                  <Input {...field} />
                                </FormControl>
                                <FormDescription>
                                  {t(`platforms.${key}`)} {t('form.descriptions.downloadLink')}
                                </FormDescription>
                              </FormItem>
                            )}
                          />
                        ))}
                      </div>
                    </div>
                  </TabsContent>
                </Tabs>
              </form>
            </Form>
          </ScrollArea>
          <SheetFooter className='flex-row justify-end gap-2 pt-3'>
            <Button variant='outline' onClick={() => setOpen(false)}>
              {t('actions.cancel')}
            </Button>
            <Button onClick={form.handleSubmit(onSubmit)} disabled={loading}>
              {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
              {editingClient ? t('actions.update') : t('actions.add')}
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>
    </>
  );
}
