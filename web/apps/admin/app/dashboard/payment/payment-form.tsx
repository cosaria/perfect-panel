'use client';

import useGlobalStore from '@/config/use-global';
import { getPaymentPlatform } from '@/services/admin/payment';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@workspace/ui/components/button';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@workspace/ui/components/form';
import { RadioGroup, RadioGroupItem } from '@workspace/ui/components/radio-group';
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
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { MarkdownEditor } from '@workspace/ui/custom-components/editor';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { unitConversion } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import * as z from 'zod';

interface PaymentFormProps<T> {
  trigger: React.ReactNode;
  title: string;
  loading?: boolean;
  initialValues?: T;
  onSubmit: (values: T) => Promise<boolean>;
  isEdit?: boolean;
}

export default function PaymentForm<T>({
  trigger,
  title,
  loading,
  initialValues,
  onSubmit,
  isEdit,
}: PaymentFormProps<T>) {
  const t = useTranslations('payment');
  const { common } = useGlobalStore();
  const { currency } = common;
  const [open, setOpen] = useState(false);

  const { data: platformData } = useQuery({
    queryKey: ['getPaymentPlatform'],
    queryFn: async () => {
      const { data } = await getPaymentPlatform();
      return data?.data?.list || [];
    },
  });

  const formSchema = z.object({
    name: z.string().min(1, { message: t('nameRequired') }),
    platform: z.string().optional(),
    icon: z.string().optional(),
    domain: z.string().optional(),
    config: z.any(),
    fee_mode: z.number().min(0).max(2),
    fee_percent: z.number().optional(),
    fee_amount: z.number().optional(),
    description: z.string().optional(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      platform: '',
      icon: '',
      domain: '',
      config: {},
      fee_mode: 0,
      fee_percent: 0,
      fee_amount: 0,
      ...(initialValues as any),
    },
  });

  const feeMode = form.watch('fee_mode');
  const platformValue = form.watch('platform');
  const configValues = form.watch('config');

  const currentPlatform = platformData?.find((p) => p.platform === platformValue);
  const currentFieldDescriptions = currentPlatform?.platform_field_description || {};
  const configFields = Object.keys(currentFieldDescriptions) || [];
  const platformUrl = currentPlatform?.platform_url || '';

  useEffect(() => {
    if (feeMode === 0) {
      form.setValue('fee_amount', 0);
      form.setValue('fee_percent', 0);
    } else if (feeMode === 1) {
      form.setValue('fee_amount', 0);
    } else if (feeMode === 2) {
      form.setValue('fee_percent', 0);
    }
  }, [feeMode, form]);

  const handleClose = () => {
    form.reset();
    setOpen(false);
  };

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    const cleanedValues = { ...values };

    if (values.fee_mode === 0) {
      cleanedValues.fee_amount = undefined;
      cleanedValues.fee_percent = undefined;
    } else if (values.fee_mode === 1) {
      cleanedValues.fee_amount = undefined;
    } else if (values.fee_mode === 2) {
      cleanedValues.fee_percent = undefined;
    }

    const success = await onSubmit(cleanedValues as unknown as T);
    if (success) {
      handleClose();
    }
  };

  const openPlatformUrl = () => {
    if (platformUrl) {
      window.open(platformUrl, '_blank');
    }
  };

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>{trigger}</SheetTrigger>
      <SheetContent className='w-[550px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{title}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100vh-48px-36px-36px-env(safe-area-inset-top))]'>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className='space-y-6 px-6 pt-4'>
              <div className='space-y-4'>
                <div className='grid grid-cols-1 gap-4 sm:grid-cols-2'>
                  <FormField
                    control={form.control}
                    name='name'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('name')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('namePlaceholder')}
                            value={field.value}
                            onValueChange={(value) => form.setValue('name', value as string)}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name='icon'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('icon')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('iconPlaceholder')}
                            value={field.value}
                            onValueChange={(value) => form.setValue('icon', value as string)}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <FormField
                  control={form.control}
                  name='domain'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('domain')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          placeholder='http(s)://example.com'
                          value={field.value}
                          onValueChange={(value) => form.setValue('domain', value as string)}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className='space-y-4'>
                <FormField
                  control={form.control}
                  name='fee_mode'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('handlingFee')}</FormLabel>
                      <FormControl>
                        <RadioGroup
                          onValueChange={(value) => field.onChange(parseInt(value))}
                          value={field.value.toString()}
                          className='flex flex-wrap gap-4'
                        >
                          <FormItem className='flex items-center space-x-2'>
                            <FormControl>
                              <RadioGroupItem value='0' />
                            </FormControl>
                            <FormLabel className='!mt-0 cursor-pointer'>{t('noFee')}</FormLabel>
                          </FormItem>
                          <FormItem className='flex items-center space-x-2'>
                            <FormControl>
                              <RadioGroupItem value='1' />
                            </FormControl>
                            <FormLabel className='!mt-0 cursor-pointer'>
                              {t('percentFee')}
                            </FormLabel>
                          </FormItem>
                          <FormItem className='flex items-center space-x-2'>
                            <FormControl>
                              <RadioGroupItem value='2' />
                            </FormControl>
                            <FormLabel className='!mt-0 cursor-pointer'>{t('fixedFee')}</FormLabel>
                          </FormItem>
                        </RadioGroup>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {feeMode === 1 && (
                  <div className='grid grid-cols-1 sm:w-1/2'>
                    <FormField
                      control={form.control}
                      name='fee_percent'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('feePercent')}</FormLabel>
                          <FormControl>
                            <EnhancedInput
                              type='number'
                              step='0.01'
                              suffix='%'
                              value={field.value}
                              onValueChange={field.onChange}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                )}

                {feeMode === 2 && (
                  <div className='grid grid-cols-1 sm:w-1/2'>
                    <FormField
                      control={form.control}
                      name='fee_amount'
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('feeAmount')}</FormLabel>
                          <FormControl>
                            <EnhancedInput
                              type='number'
                              step='0.01'
                              prefix={currency.currency_symbol}
                              suffix={currency.currency_unit}
                              value={unitConversion('centsToDollars', field.value)}
                              onValueChange={(value) =>
                                field.onChange(unitConversion('dollarsToCents', value))
                              }
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                )}
              </div>

              <div className='space-y-4'>
                {(!platformValue || platformData?.find((p) => p.platform === platformValue)) && (
                  <FormField
                    control={form.control}
                    name='platform'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('platform')}</FormLabel>
                        <Select
                          onValueChange={(value) => {
                            form.setValue('platform', value as string);
                            form.setValue('config', {});
                          }}
                          defaultValue={field.value}
                          value={field.value}
                          // @ts-expect-error - disabled prop type mismatch with SelectTrigger component
                          disabled={isEdit && Boolean(initialValues?.platform)}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder={t('selectPlatform')} />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {platformData?.map((platform) => (
                              <SelectItem key={platform.platform} value={platform.platform}>
                                {platform.platform}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        {platformUrl ? (
                          <div className='mt-1 flex justify-end'>
                            <Button
                              variant='ghost'
                              size='sm'
                              className='h-6 px-2 text-xs'
                              onClick={openPlatformUrl}
                            >
                              <Icon icon='tabler:external-link' className='mr-1 h-3 w-3' />
                              {t('applyForPayment')}
                            </Button>
                          </div>
                        ) : (
                          <div className='mt-1 h-6'></div>
                        )}
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                )}

                {configFields.length > 0 && (
                  <div className='mt-4 space-y-4'>
                    {configFields.map((fieldKey) => (
                      <FormItem key={fieldKey}>
                        <FormLabel>{currentFieldDescriptions[fieldKey]}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('configPlaceholder', {
                              field: currentFieldDescriptions[fieldKey],
                            })}
                            value={
                              configValues && configValues[fieldKey] !== undefined
                                ? configValues[fieldKey]
                                : ''
                            }
                            disabled={fieldKey === 'webhook_secret'}
                            onValueChange={(value) => {
                              const newConfig = { ...configValues };
                              newConfig[fieldKey] = value;
                              form.setValue('config', newConfig);
                            }}
                          />
                        </FormControl>
                      </FormItem>
                    ))}
                  </div>
                )}

                <FormField
                  control={form.control}
                  name='description'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('description')}</FormLabel>
                      <FormControl>
                        <MarkdownEditor
                          value={field.value}
                          onChange={(value) => form.setValue(field.name, value as string)}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </form>
          </Form>
        </ScrollArea>

        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={handleClose}>
            {t('cancel')}
          </Button>
          <Button disabled={loading} onClick={form.handleSubmit(handleSubmit)}>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('submit')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
