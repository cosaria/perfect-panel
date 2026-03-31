'use client';

import { getVerifyConfig, updateVerifyConfig } from '@/services/admin/system';
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
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const verifySchema = z.object({
  turnstile_site_key: z.string().optional(),
  turnstile_secret: z.string().optional(),
  enable_register_verify: z.boolean().optional(),
  enable_login_verify: z.boolean().optional(),
  enable_password_verify: z.boolean().optional(),
});

type VerifyFormData = z.infer<typeof verifySchema>;

export default function VerifyConfig() {
  const t = useTranslations('system.verify');
  const systemT = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getVerifyConfig'],
    queryFn: async () => {
      const { data } = await getVerifyConfig();
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<VerifyFormData>({
    resolver: zodResolver(verifySchema),
    defaultValues: {
      turnstile_site_key: '',
      turnstile_secret: '',
      enable_register_verify: false,
      enable_login_verify: false,
      enable_password_verify: false,
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: VerifyFormData) {
    setLoading(true);
    try {
      await updateVerifyConfig(values as API.VerifyConfig);
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
              <Icon icon='mdi:shield-check-outline' className='text-primary h-5 w-5' />
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
              id='verify-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='turnstile_site_key'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('turnstileSiteKey')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('turnstileSiteKeyPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('turnstileSiteKeyDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='turnstile_secret'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('turnstileSecret')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('turnstileSecretPlaceholder')}
                        value={field.value}
                        type='password'
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('turnstileSecretDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='enable_register_verify'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('enableRegisterVerify')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('enableRegisterVerifyDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='enable_login_verify'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('enableLoginVerify')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('enableLoginVerifyDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='enable_password_verify'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('enablePasswordVerify')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('enablePasswordVerifyDescription')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='verify-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {systemT('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
