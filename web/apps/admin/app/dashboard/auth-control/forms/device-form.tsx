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
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { uid } from 'radash';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const deviceSchema = z.object({
  id: z.number(),
  method: z.string(),
  enabled: z.boolean(),
  config: z
    .object({
      show_ads: z.boolean().optional(),
      only_real_device: z.boolean().optional(),
      enable_security: z.boolean().optional(),
      security_secret: z.string().optional(),
    })
    .optional(),
});

type DeviceFormData = z.infer<typeof deviceSchema>;

export default function DeviceForm() {
  const t = useTranslations('auth-control');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getAuthMethodConfig', 'device'],
    queryFn: async () => {
      const { data } = await getAuthMethodConfig({
        method: 'device',
      });
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<DeviceFormData>({
    resolver: zodResolver(deviceSchema),
    defaultValues: {
      id: 0,
      method: 'device',
      enabled: false,
      config: {
        show_ads: false,
        only_real_device: false,
        enable_security: false,
        security_secret: '',
      },
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: DeviceFormData) {
    setLoading(true);
    try {
      await updateAuthMethodConfig(values as API.UpdateAuthMethodConfigRequest);
      toast.success(t('common.saveSuccess'));
      refetch();
      setOpen(false);
    } catch (error) {
      toast.error(t('common.saveFailed'));
    } finally {
      setLoading(false);
    }
  }

  function generateSecurityKey() {
    const id = uid(32).toLowerCase();
    const formatted = `${id.slice(0, 8)}-${id.slice(8, 12)}-${id.slice(12, 16)}-${id.slice(16, 20)}-${id.slice(20)}`;
    form.setValue('config.security_secret', formatted);
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <div className='flex cursor-pointer items-center justify-between transition-colors'>
          <div className='flex items-center gap-3'>
            <div className='bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg'>
              <Icon icon='mdi:devices' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('device.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('device.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('device.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='device-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='enabled'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('device.enable')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('device.enableDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.show_ads'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('device.showAds')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('device.showAdsDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.only_real_device'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('device.blockVirtualMachine')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('device.blockVirtualMachineDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.enable_security'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('device.enableSecurity')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('device.enableSecurityDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.security_secret'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('device.communicationKey')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='e.g., 12345678-1234-1234-1234-123456789abc'
                        value={field.value}
                        onValueChange={field.onChange}
                        suffix={
                          <div className='bg-muted flex h-9 items-center text-nowrap px-3'>
                            <Icon
                              icon='mdi:dice-multiple'
                              onClick={generateSecurityKey}
                              className='size-4 cursor-pointer'
                            />
                          </div>
                        }
                      />
                    </FormControl>
                    <FormDescription>{t('device.communicationKeyDescription')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='device-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
