'use client';

import { getRegisterConfig, updateRegisterConfig } from '@/services/admin/system';
import { useSubscribe } from '@/store/subscribe';
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
import { Combobox } from '@workspace/ui/custom-components/combobox';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const registerSchema = z.object({
  stop_register: z.boolean().optional(),
  enable_trial: z.boolean().optional(),
  trial_subscribe: z.number().optional(),
  trial_time: z.number().optional(),
  trial_time_unit: z.string().optional(),
  enable_ip_register_limit: z.boolean().optional(),
  ip_register_limit: z.number().optional(),
  ip_register_limit_duration: z.number().optional(),
});

type RegisterFormData = z.infer<typeof registerSchema>;

export default function RegisterConfig() {
  const t = useTranslations('system.register');
  const systemT = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getRegisterConfig'],
    queryFn: async () => {
      const { data } = await getRegisterConfig();
      return data.data;
    },
    enabled: open,
  });

  const { subscribes } = useSubscribe();

  const form = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      stop_register: false,
      enable_trial: false,
      trial_subscribe: undefined,
      trial_time: 0,
      trial_time_unit: 'day',
      enable_ip_register_limit: false,
      ip_register_limit: 1,
      ip_register_limit_duration: 1,
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: RegisterFormData) {
    setLoading(true);
    try {
      await updateRegisterConfig(values as API.RegisterConfig);
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
              <Icon icon='mdi:account-plus-outline' className='text-primary h-5 w-5' />
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
              id='register-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='stop_register'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('stopNewUserRegistration')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('stopNewUserRegistrationDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='enable_ip_register_limit'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('ipRegistrationLimit')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('ipRegistrationLimitDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {form.watch('enable_ip_register_limit') && (
                <div className='grid grid-cols-1 gap-4 md:grid-cols-2'>
                  <FormField
                    control={form.control}
                    name='ip_register_limit'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('registrationLimitCount')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('inputPlaceholder')}
                            value={field.value}
                            type='number'
                            min={1}
                            onValueBlur={(value) => field.onChange(Number(value))}
                          />
                        </FormControl>
                        <FormDescription>{t('registrationLimitCountDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='ip_register_limit_duration'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('registrationLimitExpire')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('inputPlaceholder')}
                            value={field.value}
                            type='number'
                            min={1}
                            suffix={t('minute')}
                            onValueBlur={(value) => field.onChange(Number(value))}
                          />
                        </FormControl>
                        <FormDescription>{t('registrationLimitExpireDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              )}

              <FormField
                control={form.control}
                name='enable_trial'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('enableTrial')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('enableTrialDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {form.watch('enable_trial') && (
                <>
                  <FormField
                    control={form.control}
                    name='trial_time'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('trialConfig')}</FormLabel>
                        <FormControl>
                          <div className='flex gap-2'>
                            <EnhancedInput
                              placeholder={t('inputPlaceholder')}
                              value={field.value}
                              type='number'
                              min={0}
                              onValueBlur={(value) => field.onChange(Number(value))}
                              className='flex-1'
                              prefix={
                                <FormField
                                  control={form.control}
                                  name='trial_subscribe'
                                  render={({ field }) => (
                                    <Combobox
                                      placeholder={t('selectPlaceholder')}
                                      value={field.value}
                                      onChange={(value: number) => {
                                        if (value) {
                                          field.onChange(value);
                                        }
                                      }}
                                      options={subscribes?.map((item) => ({
                                        label: item.name!,
                                        value: item.id!,
                                      }))}
                                      className='bg-secondary w-32 rounded-r-none'
                                    />
                                  )}
                                />
                              }
                              suffix={
                                <FormField
                                  control={form.control}
                                  name='trial_time_unit'
                                  render={({ field: unitField }) => (
                                    <Combobox
                                      placeholder={t('selectPlaceholder')}
                                      value={unitField.value}
                                      onChange={(value: string) => {
                                        unitField.onChange(value);
                                      }}
                                      options={[
                                        { label: t('none'), value: 'None' },
                                        { label: t('year'), value: 'Year' },
                                        { label: t('month'), value: 'Month' },
                                        { label: t('day'), value: 'Day' },
                                        { label: t('hour'), value: 'Hour' },
                                        { label: t('minute'), value: 'Minute' },
                                      ]}
                                      className='bg-secondary w-32 rounded-l-none'
                                    />
                                  )}
                                />
                              }
                            />
                          </div>
                        </FormControl>
                        <FormDescription>{t('trialConfigDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </>
              )}
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {systemT('common.cancel')}
          </Button>
          <Button disabled={loading} type='submit' form='register-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {systemT('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
