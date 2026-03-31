'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

import { getSubscribeConfig, updateSubscribeConfig } from '@/services/admin/system';
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

const subscribeConfigSchema = z.object({
  single_model: z.boolean().optional(),
  pan_domain: z.boolean().optional(),
  subscribe_path: z.string().optional(),
  subscribe_domain: z.string().optional(),
  user_agent_limit: z.boolean().optional(),
  user_agent_list: z.string().optional(),
});

type SubscribeConfigFormData = z.infer<typeof subscribeConfigSchema>;

export default function ConfigForm() {
  const t = useTranslations('subscribe');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getSubscribeConfig'],
    queryFn: async () => {
      const { data } = await getSubscribeConfig();
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<SubscribeConfigFormData>({
    resolver: zodResolver(subscribeConfigSchema),
    defaultValues: {
      single_model: false,
      pan_domain: false,
      subscribe_path: '',
      subscribe_domain: '',
      user_agent_limit: false,
      user_agent_list: '',
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: SubscribeConfigFormData) {
    setLoading(true);
    try {
      await updateSubscribeConfig(values as API.SubscribeConfig);
      toast.success(t('config.updateSuccess'));
      refetch();
      setOpen(false);
    } catch (error) {
      toast.error(t('config.updateError'));
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
              <Icon icon='mdi:cog' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('config.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('config.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('config.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='subscribe-config-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='single_model'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('config.singleSubscriptionMode')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>
                      {t('config.singleSubscriptionModeDescription')}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='pan_domain'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('config.wildcardResolution')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('config.wildcardResolutionDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='subscribe_path'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('config.subscriptionPath')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('config.subscriptionPathPlaceholder')}
                        value={field.value}
                        onValueBlur={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('config.subscriptionPathDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='subscribe_domain'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('config.subscriptionDomain')}</FormLabel>
                    <FormControl>
                      <Textarea
                        className='h-32'
                        placeholder={`${t('config.subscriptionDomainPlaceholder')}\nexample.com\nwww.example.com`}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>{t('config.subscriptionDomainDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='user_agent_limit'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('config.userAgentLimit', { userAgent: 'User-Agent' })}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>
                      {t('config.userAgentLimitDescription', { userAgent: 'User-Agent' })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='user_agent_list'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t('config.userAgentList', {
                        userAgent: 'User-Agent',
                      })}
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        className='h-32'
                        placeholder={`${t('config.userAgentListPlaceholder', { userAgent: 'User-Agent' })}\nClashX\nClashForAndroid\nClash-verge`}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t('config.userAgentListDescription', { userAgent: 'User-Agent' })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {t('actions.cancel')}
          </Button>
          <Button disabled={loading} type='submit' form='subscribe-config-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('actions.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
