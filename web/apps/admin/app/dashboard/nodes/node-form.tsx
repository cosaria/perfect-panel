'use client';

import { useNode } from '@/store/node';
import { useServer } from '@/store/server';
import { zodResolver } from '@hookform/resolvers/zod';
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
import { ScrollArea } from '@workspace/ui/components/scroll-area';
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { Combobox } from '@workspace/ui/custom-components/combobox';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import TagInput from '@workspace/ui/custom-components/tag-input';
import { useTranslations } from 'next-intl';
import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

export type ProtocolName =
  | 'shadowsocks'
  | 'vmess'
  | 'vless'
  | 'trojan'
  | 'hysteria'
  | 'tuic'
  | 'anytls'
  | 'naive'
  | 'http'
  | 'socks'
  | 'mieru';

const buildSchema = (t: ReturnType<typeof useTranslations>) =>
  z.object({
    name: z.string().trim().min(1, t('errors.nameRequired')),
    server_id: z
      .number({ message: t('errors.serverRequired') })
      .int()
      .gt(0, t('errors.serverRequired'))
      .optional(),
    protocol: z.string().min(1, t('errors.protocolRequired')),
    address: z.string().trim().min(1, t('errors.serverAddrRequired')),
    port: z
      .number({ message: t('errors.portRange') })
      .int()
      .min(1, t('errors.portRange'))
      .max(65535, t('errors.portRange')),
    tags: z.array(z.string()),
  });

export type NodeFormValues = z.infer<ReturnType<typeof buildSchema>>;

export default function NodeForm(props: {
  trigger: string;
  title: string;
  loading?: boolean;
  initialValues?: Partial<NodeFormValues>;
  onSubmit: (values: NodeFormValues) => Promise<boolean> | boolean;
}) {
  const { trigger, title, loading, initialValues, onSubmit } = props;
  const t = useTranslations('nodes');
  const Scheme = useMemo(() => buildSchema(t), [t]);
  const [open, setOpen] = useState(false);

  const [autoFilledFields, setAutoFilledFields] = useState<Set<string>>(new Set());

  const addAutoFilledField = (fieldName: string) => {
    setAutoFilledFields((prev) => new Set(prev).add(fieldName));
  };

  const removeAutoFilledField = (fieldName: string) => {
    setAutoFilledFields((prev) => {
      const newSet = new Set(prev);
      newSet.delete(fieldName);
      return newSet;
    });
  };

  const form = useForm<NodeFormValues>({
    resolver: zodResolver(Scheme),
    defaultValues: {
      name: '',
      server_id: undefined,
      protocol: '',
      address: '',
      port: 0,
      tags: [],
      ...initialValues,
    },
  });

  const serverId = form.watch('server_id');

  const { servers, getAvailableProtocols } = useServer();
  const { tags } = useNode();

  const existingTags: string[] = tags || [];

  const availableProtocols = getAvailableProtocols(serverId);

  useEffect(() => {
    if (initialValues) {
      form.reset({
        name: '',
        server_id: undefined,
        protocol: '',
        address: '',
        port: 0,
        tags: [],
        ...initialValues,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialValues]);

  function handleServerChange(nextId?: number | null) {
    const id = nextId ?? undefined;
    form.setValue('server_id', id);

    if (!id) {
      setAutoFilledFields(new Set());
      return;
    }

    const selectedServer = servers.find((s) => s.id === id);
    if (!selectedServer) return;

    const currentValues = form.getValues();
    const fieldsToFill: string[] = [];

    if (!currentValues.name || autoFilledFields.has('name')) {
      form.setValue('name', selectedServer.name as string, { shouldDirty: false });
      fieldsToFill.push('name');
    }

    if (!currentValues.address || autoFilledFields.has('address')) {
      form.setValue('address', selectedServer.address as string, { shouldDirty: false });
      fieldsToFill.push('address');
    }

    const protocols = getAvailableProtocols(id);
    const firstProtocol = protocols[0];

    if (firstProtocol && (!currentValues.protocol || autoFilledFields.has('protocol'))) {
      form.setValue('protocol', firstProtocol.protocol, { shouldDirty: false });
      fieldsToFill.push('protocol');

      if (!currentValues.port || currentValues.port === 0 || autoFilledFields.has('port')) {
        const port = firstProtocol.port || 0;
        form.setValue('port', port, { shouldDirty: false });
        fieldsToFill.push('port');
      }
    }

    setAutoFilledFields(new Set(fieldsToFill));
  }

  const handleManualFieldChange = (fieldName: keyof NodeFormValues, value: any) => {
    form.setValue(fieldName, value);
    removeAutoFilledField(fieldName);
  };

  function handleProtocolChange(nextProto?: ProtocolName | null) {
    const protocol = (nextProto || '') as ProtocolName | '';
    form.setValue('protocol', protocol);

    if (!protocol || !serverId) {
      removeAutoFilledField('protocol');
      return;
    }

    const currentValues = form.getValues();
    const isPortAutoFilled = autoFilledFields.has('port');

    removeAutoFilledField('protocol');

    if (!currentValues.port || currentValues.port === 0 || isPortAutoFilled) {
      const protocolData = availableProtocols.find((p) => p.protocol === protocol);

      if (protocolData) {
        const port = protocolData.port || 0;
        form.setValue('port', port, { shouldDirty: false });
        addAutoFilledField('port');
      }
    }
  }

  async function handleSubmit(values: NodeFormValues) {
    const result = await onSubmit(values);
    if (result) {
      setOpen(false);
      setAutoFilledFields(new Set());
    }
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          onClick={() => {
            form.reset();
            setAutoFilledFields(new Set());
          }}
        >
          {trigger}
        </Button>
      </SheetTrigger>

      <SheetContent className='w-[560px] max-w-full'>
        <SheetHeader>
          <SheetTitle>{title}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6 pt-4'>
          <Form {...form}>
            <form className='grid grid-cols-1 gap-4'>
              <FormField
                control={form.control}
                name='server_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('server')}</FormLabel>
                    <FormControl>
                      <Combobox<number, false>
                        placeholder={t('select_server')}
                        value={field.value}
                        options={servers.map((s) => ({
                          value: s.id,
                          label: `${s.name} (${(s.address as any) || ''})`,
                        }))}
                        onChange={(v) => handleServerChange(v)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='protocol'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('protocol')}</FormLabel>
                    <FormControl>
                      <Combobox<string, false>
                        placeholder={t('select_protocol')}
                        value={field.value}
                        options={availableProtocols.map((p) => ({
                          value: p.protocol,
                          label: `${p.protocol}${p.port ? ` (${p.port})` : ''}`,
                        }))}
                        onChange={(v) => handleProtocolChange((v as ProtocolName) || null)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('name')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        {...field}
                        onValueChange={(v) => handleManualFieldChange('name', v as string)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='address'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('address')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        {...field}
                        onValueChange={(v) => handleManualFieldChange('address', v as string)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='port'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('port')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        {...field}
                        type='number'
                        min={1}
                        max={65535}
                        placeholder='1-65535'
                        onValueChange={(v) => handleManualFieldChange('port', Number(v))}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='tags'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('tags')}</FormLabel>
                    <FormControl>
                      <TagInput
                        placeholder={t('tags_placeholder')}
                        value={field.value || []}
                        onChange={(v) => form.setValue(field.name, v)}
                        options={existingTags}
                      />
                    </FormControl>
                    <FormDescription>{t('tags_description')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </ScrollArea>

        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {t('cancel')}
          </Button>
          <Button
            disabled={loading}
            onClick={form.handleSubmit(handleSubmit, (errors) => {
              const key = Object.keys(errors)[0] as keyof typeof errors;
              if (key) toast.error(String(errors[key]?.message));
              return false;
            })}
          >
            {t('confirm')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
