'use client';

import { getVerifyCodeConfig, updateVerifyCodeConfig } from '@/services/admin/system';
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
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const verifyCodeSchema = z.object({
  verify_code_expire_time: z.number().optional(),
  verify_code_interval: z.number().optional(),
  verify_code_limit: z.number().optional(),
});

type VerifyCodeFormData = z.infer<typeof verifyCodeSchema>;

export default function VerifyCodeConfig() {
  const t = useTranslations('system.verifyCode');
  const systemT = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getVerifyCodeConfig'],
    queryFn: async () => {
      const { data } = await getVerifyCodeConfig();
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<VerifyCodeFormData>({
    resolver: zodResolver(verifyCodeSchema),
    defaultValues: {
      verify_code_expire_time: 300,
      verify_code_interval: 60,
      verify_code_limit: 10,
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: VerifyCodeFormData) {
    setLoading(true);
    try {
      await updateVerifyCodeConfig(values as API.VerifyCodeConfig);
      toast.success(t('saveSuccess'));
      refetch();
      setOpen(false);
    } catch (error) {
      toast.error(t('saveFailed'));
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
              <Icon icon='mdi:message-text-clock-outline' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('title')}</p>
              <p className='text-muted-foreground text-sm'>{t('description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='verify-code-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='verify_code_expire_time'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('expireTime')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('inputPlaceholder')}
                        value={field.value}
                        type='number'
                        min={60}
                        suffix={t('seconds')}
                        onValueBlur={(value) => field.onChange(Number(value))}
                      />
                    </FormControl>
                    <FormDescription>{t('expireTimeDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='verify_code_interval'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('interval')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('inputPlaceholder')}
                        value={field.value}
                        type='number'
                        min={30}
                        suffix={t('seconds')}
                        onValueBlur={(value) => field.onChange(Number(value))}
                      />
                    </FormControl>
                    <FormDescription>{t('intervalDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='verify_code_limit'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('dailyLimit')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('inputPlaceholder')}
                        value={field.value}
                        type='number'
                        min={1}
                        suffix={t('times')}
                        onValueBlur={(value) => field.onChange(Number(value))}
                      />
                    </FormControl>
                    <FormDescription>{t('dailyLimitDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {systemT('common.cancel')}
          </Button>
          <Button disabled={loading} type='submit' form='verify-code-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {systemT('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
