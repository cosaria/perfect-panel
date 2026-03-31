'use client';

import useGlobalStore from '@/config/use-global';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@workspace/ui/components/button';
import {
  Form,
  FormControl,
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
import { AreaCodeSelect } from '@workspace/ui/custom-components/area-code-select';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { unitConversion } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

interface UserFormProps<T> {
  onSubmit: (data: T) => Promise<boolean> | boolean;
  initialValues?: T;
  loading?: boolean;
  trigger: string;
  title: string;
  update?: boolean;
}

export default function UserForm<T extends Record<string, any>>({
  onSubmit,
  initialValues,
  loading,
  trigger,
  title,
}: Readonly<UserFormProps<T>>) {
  const t = useTranslations('user');
  const { common } = useGlobalStore();
  const { currency } = common;

  const [open, setOpen] = useState(false);
  const formSchema = z.object({
    email: z.string().email(t('invalidEmailFormat')),
    telephone_area_code: z.string().optional(),
    telephone: z.string().optional(),
    password: z.string().optional(),
    referer_id: z.number().optional(),
    refer_code: z.string().optional(),
    referral_percentage: z.number().optional(),
    only_first_purchase: z.boolean().optional(),
    is_admin: z.boolean().optional(),
    balance: z.number().optional(),
    gift_amount: z.number().optional(),
    commission: z.number().optional(),
  });
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...initialValues,
    },
  });

  useEffect(() => {
    form?.reset(initialValues);
  }, [form, initialValues]);

  async function handleSubmit(data: { [x: string]: any }) {
    const bool = await onSubmit(data as T);

    if (bool) setOpen(false);
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          onClick={() => {
            form.reset();
            setOpen(true);
          }}
        >
          {trigger}
        </Button>
      </SheetTrigger>
      <SheetContent className='w-[500px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{title}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))]'>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className='space-y-4 px-6 pt-4'>
              <FormField
                control={form.control}
                name='email'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userEmail')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('userEmailPlaceholder')}
                        {...field}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='telephone'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('telephone')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        prefix={
                          <FormField
                            control={form.control}
                            name='telephone_area_code'
                            render={({ field }) => (
                              <FormItem>
                                <FormControl>
                                  <AreaCodeSelect
                                    className='w-32 rounded-none border-y-0 border-l-0'
                                    simple
                                    placeholder={t('areaCodePlaceholder')}
                                    value={field.value}
                                    onChange={(value) => {
                                      form.setValue(field.name, value.phone as string);
                                    }}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        }
                        placeholder={t('telephonePlaceholder')}
                        {...field}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='password'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('password')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        autoComplete='new-password'
                        placeholder={t('passwordPlaceholder')}
                        {...field}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='referer_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('refererId')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('refererIdPlaceholder')}
                        {...field}
                        type='number'
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='refer_code'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('inviteCode')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('inviteCodePlaceholder')}
                        {...field}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='referral_percentage'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('referralPercentage')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='number'
                        min={0}
                        max={100}
                        placeholder={t('referralPercentagePlaceholder')}
                        {...field}
                        suffix='%'
                        onValueChange={(value) => {
                          form.setValue(field.name, Number(value));
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='only_first_purchase'
                render={({ field }) => (
                  <FormItem className='flex items-center justify-between space-x-2'>
                    <FormLabel>{t('onlyFirstPurchase')}</FormLabel>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='balance'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('balance')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        prefix={currency?.currency_symbol ?? '$'}
                        placeholder={t('balancePlaceholder')}
                        type='number'
                        {...field}
                        formatInput={(value) => unitConversion('centsToDollars', value)}
                        formatOutput={(value) => unitConversion('dollarsToCents', value)}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                        min={0}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='gift_amount'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('giftAmount')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        prefix={currency?.currency_symbol ?? '$'}
                        placeholder={t('giftAmountPlaceholder')}
                        type='number'
                        {...field}
                        formatInput={(value) => unitConversion('centsToDollars', value)}
                        formatOutput={(value) => unitConversion('dollarsToCents', value)}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                        min={0}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='commission'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('commission')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        prefix={currency?.currency_symbol ?? '$'}
                        placeholder={t('commissionPlaceholder')}
                        type='number'
                        {...field}
                        formatInput={(value) => unitConversion('centsToDollars', value)}
                        formatOutput={(value) => unitConversion('dollarsToCents', value)}
                        onValueChange={(value) => {
                          form.setValue(field.name, value);
                        }}
                        min={0}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='is_admin'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('manager')}</FormLabel>
                    <FormControl>
                      <div className='pt-2'>
                        <Switch checked={!!field.value} onCheckedChange={field.onChange} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button
            variant='outline'
            disabled={loading}
            onClick={() => {
              setOpen(false);
            }}
          >
            {t('cancel')}
          </Button>
          <Button disabled={loading} onClick={form.handleSubmit(handleSubmit)}>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />} {t('confirm')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
