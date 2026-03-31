'use client';

import { getAuthMethodConfig, updateAuthMethodConfig } from '@/services/admin/authMethod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
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
import { Switch } from '@workspace/ui/components/switch';
import { Textarea } from '@workspace/ui/components/textarea';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const appleSchema = z.object({
  enabled: z.boolean(),
  config: z
    .object({
      team_id: z.string().optional(),
      key_id: z.string().optional(),
      client_id: z.string().optional(),
      client_secret: z.string().optional(),
      redirect_url: z.string().optional(),
    })
    .optional(),
});

type AppleFormData = z.infer<typeof appleSchema>;

export default function AppleForm() {
  const t = useTranslations('auth-control');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getAuthMethodConfig', 'apple'],
    queryFn: async () => {
      const { data } = await getAuthMethodConfig({
        method: 'apple',
      });

      return data.data;
    },
    enabled: open,
  });

  const form = useForm<AppleFormData>({
    resolver: zodResolver(appleSchema),
    defaultValues: {
      enabled: false,
      config: {
        team_id: '',
        key_id: '',
        client_id: '',
        client_secret: '',
        redirect_url: '',
      },
    },
  });

  useEffect(() => {
    if (data) {
      form.reset({
        enabled: data.enabled || false,
        config: {
          team_id: data.config?.team_id || '',
          key_id: data.config?.key_id || '',
          client_id: data.config?.client_id || '',
          client_secret: data.config?.client_secret || '',
          redirect_url: data.config?.redirect_url || '',
        },
      });
    }
  }, [data, form]);

  async function onSubmit(values: AppleFormData) {
    setLoading(true);
    try {
      await updateAuthMethodConfig({
        ...data,
        enabled: values.enabled,
        config: {
          ...data?.config,
          ...values.config,
        },
      } as API.UpdateAuthMethodConfigRequest);
      toast.success(t('common.saveSuccess'));
      refetch();
      setOpen(false);
    } catch (error) {
      toast.error(t('common.saveFailed'));
    } finally {
      setLoading(false);
    }
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <div className='flex cursor-pointer items-center justify-between transition-colors'>
          <div className='flex items-center gap-3'>
            <div className='bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg'>
              <Icon icon='mdi:apple' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('apple.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('apple.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[500px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('apple.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form id='apple-form' onSubmit={form.handleSubmit(onSubmit)} className='space-y-2 pt-4'>
              <FormField
                control={form.control}
                name='enabled'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('apple.enable')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('apple.enableDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.team_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('apple.teamId')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='ABCDE1FGHI'
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('apple.teamIdDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.key_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('apple.keyId')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='ABC1234567'
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('apple.keyIdDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.client_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('apple.clientId')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='com.your.app.service'
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('apple.clientIdDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.client_secret'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('apple.clientSecret')}</FormLabel>
                    <FormControl>
                      <Textarea
                        className='h-20'
                        placeholder={`-----BEGIN PRIVATE KEY-----\nMIGTAgEA...\n-----END PRIVATE KEY-----`}
                        value={field.value}
                        onChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('apple.clientSecretDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.redirect_url'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('apple.redirectUri')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='https://your-domain.com'
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('apple.redirectUriDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {t('common.cancel')}
          </Button>
          <Button disabled={loading} type='submit' form='apple-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
