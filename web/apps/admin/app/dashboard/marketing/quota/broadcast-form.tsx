'use client';

import { Display } from '@/components/display';
import { createQuotaTask, queryQuotaTaskPreCount } from '@/services/admin/marketing';
import { useSubscribe } from '@/store/subscribe';
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
import { RadioGroup, RadioGroupItem } from '@workspace/ui/components/radio-group';
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
import { unitConversion } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

export default function QuotaBroadcastForm() {
  const t = useTranslations('marketing');

  // Define schema with internationalized error messages
  const quotaBroadcastSchema = z.object({
    subscribers: z.array(z.number()).min(1, t('pleaseSelectSubscribers')),
    is_active: z.boolean(),
    start_time: z.string().optional(),
    end_time: z.string().optional(),
    reset_traffic: z.boolean(),
    days: z.number().optional(),
    gift_type: z.number(),
    gift_value: z.number().optional(),
  });

  type QuotaBroadcastFormData = z.infer<typeof quotaBroadcastSchema>;

  const form = useForm<QuotaBroadcastFormData>({
    resolver: zodResolver(quotaBroadcastSchema),
    mode: 'onChange', // Enable real-time validation
    defaultValues: {
      subscribers: [],
      is_active: true,
      start_time: '',
      end_time: '',
      reset_traffic: false,
      days: 0,
      gift_type: 1,
      gift_value: 0,
    },
  });

  const [recipients, setRecipients] = useState<number>(0);
  const [isCalculating, setIsCalculating] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [open, setOpen] = useState(false);

  const { subscribes } = useSubscribe();

  // Calculate recipient count
  const calculateRecipients = async () => {
    setIsCalculating(true);
    try {
      const formData = form.getValues();
      let start_time: number = 0;
      let end_time: number = 0;

      if (formData.start_time) {
        start_time = new Date(formData.start_time).getTime();
      }

      if (formData.end_time) {
        end_time = new Date(formData.end_time).getTime();
      }

      const response = await queryQuotaTaskPreCount({
        subscribers: formData.subscribers,
        is_active: formData.is_active,
        start_time,
        end_time,
      });

      if (response.data?.data?.count !== undefined) {
        setRecipients(response.data.data.count);
      }
    } catch (error) {
      console.error('Failed to calculate recipients:', error);
      toast.error(t('failedToCalculateRecipients'));
      setRecipients(0);
    } finally {
      setIsCalculating(false);
    }
  };

  // Watch form values and recalculate recipients only when sheet is open
  const watchedValues = form.watch();

  useEffect(() => {
    if (!open) return; // Only calculate when sheet is open

    const debounceTimer = setTimeout(() => {
      calculateRecipients();
    }, 500); // Add debounce to avoid too frequent API calls

    return () => clearTimeout(debounceTimer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    open,
    watchedValues.subscribers,
    watchedValues.is_active,
    watchedValues.start_time,
    watchedValues.end_time,
  ]);

  const onSubmit = async (data: QuotaBroadcastFormData) => {
    setIsSubmitting(true);
    try {
      let start_time: number = 0;
      let end_time: number = 0;

      if (data.start_time) {
        start_time = Math.floor(new Date(data.start_time).getTime());
      }

      if (data.end_time) {
        end_time = Math.floor(new Date(data.end_time).getTime());
      }

      await createQuotaTask({
        subscribers: data.subscribers,
        is_active: data.is_active,
        start_time,
        end_time,
        reset_traffic: data.reset_traffic,
        days: data.days || 0,
        gift_type: data.gift_type,
        gift_value: data.gift_value || 0,
      });

      toast.success(t('quotaTaskCreatedSuccessfully'));
      form.reset();
      setRecipients(0);
      setOpen(false); // Close the sheet after successful submission
    } catch (error) {
      console.error('Failed to create quota task:', error);
      toast.error(t('failedToCreateQuotaTask'));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <div className='flex cursor-pointer items-center justify-between transition-colors'>
          <div className='flex items-center gap-3'>
            <div className='bg-primary/10 flex h-10 w-10 items-center justify-center rounded-lg'>
              <Icon icon='mdi:gift' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('quotaBroadcast')}</p>
              <p className='text-muted-foreground text-sm'>{t('createAndSendQuotaTasks')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('createQuotaTask')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-32px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='quota-broadcast-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='mt-4 space-y-6'
            >
              {/* Subscribers selection */}
              <FormField
                control={form.control}
                name='subscribers'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('subscribers')}</FormLabel>
                    <FormControl>
                      <Combobox
                        multiple={true}
                        value={field.value || []}
                        onChange={field.onChange}
                        placeholder={t('pleaseSelectSubscribers')}
                        options={subscribes?.map((subscribe) => ({
                          value: subscribe.id!,
                          label: subscribe.name!,
                          children: (
                            <div>
                              <div>{subscribe.name}</div>
                              <div className='text-muted-foreground text-xs'>
                                <Display type='traffic' value={subscribe.traffic || 0} /> /{' '}
                                <Display type='currency' value={subscribe.unit_price || 0} />
                              </div>
                            </div>
                          ),
                        }))}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Subscription count info and active status */}
              <div className='grid grid-cols-1 gap-4 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='is_active'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('validOnly')}</FormLabel>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={field.onChange}
                          className='float-end !mt-0'
                        />
                      </FormControl>
                      <FormDescription>{t('selectValidSubscriptionsOnly')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className='border-l-primary bg-primary/10 flex items-center border-l-4 px-4 py-3 text-sm'>
                  <span className='text-muted-foreground'>{t('subscriptionCount')}: </span>
                  <span className='text-primary text-lg font-medium'>
                    {isCalculating ? (
                      <Icon icon='mdi:loading' className='ml-2 h-4 w-4 animate-spin' />
                    ) : (
                      recipients.toLocaleString()
                    )}
                  </span>
                </div>
              </div>

              {/* Subscription validity period range */}
              <div className='grid grid-cols-2 gap-4'>
                <FormField
                  control={form.control}
                  name='start_time'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('subscriptionValidityStartDate')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          type='datetime-local'
                          step='1'
                          value={field.value}
                          onValueChange={field.onChange}
                        />
                      </FormControl>
                      <FormDescription>{t('includeSubscriptionsValidAfter')}</FormDescription>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='end_time'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('subscriptionValidityEndDate')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          type='datetime-local'
                          step='1'
                          value={field.value}
                          onValueChange={field.onChange}
                        />
                      </FormControl>
                      <FormDescription>{t('includeSubscriptionsValidBefore')}</FormDescription>
                    </FormItem>
                  )}
                />
              </div>

              {/* Reset traffic */}
              <FormField
                control={form.control}
                name='reset_traffic'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('resetTraffic')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('resetTrafficDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Quota days */}
              <FormField
                control={form.control}
                name='days'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('quotaDays')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='number'
                        min={1}
                        value={field.value?.toString()}
                        onValueChange={(value) => field.onChange(parseInt(value, 10))}
                      />
                    </FormControl>
                    <FormDescription>{t('numberOfDaysForTheQuota')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Gift configuration */}
              <FormField
                control={form.control}
                name='gift_type'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('giftType')}</FormLabel>
                    <FormControl>
                      <RadioGroup
                        defaultValue={String(field.value)}
                        onValueChange={(value) => {
                          field.onChange(Number(value));
                          form.setValue('gift_value', 0);
                        }}
                        className='flex gap-4'
                      >
                        <FormItem className='flex items-center space-x-3 space-y-0'>
                          <FormControl>
                            <RadioGroupItem value='1' />
                          </FormControl>
                          <FormLabel className='font-normal'>{t('fixedAmount')}</FormLabel>
                        </FormItem>
                        <FormItem className='flex items-center space-x-3 space-y-0'>
                          <FormControl>
                            <RadioGroupItem value='2' />
                          </FormControl>
                          <FormLabel className='font-normal'>{t('percentageAmount')}</FormLabel>
                        </FormItem>
                      </RadioGroup>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Gift amount based on type */}
              {form.watch('gift_type') === 1 && (
                <FormField
                  control={form.control}
                  name='gift_value'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('giftAmount')}</FormLabel>
                      <FormControl>
                        <EnhancedInput<number>
                          placeholder={t('enterAmount')}
                          type='number'
                          value={field.value}
                          formatInput={(value) => unitConversion('centsToDollars', value)}
                          formatOutput={(value) => unitConversion('dollarsToCents', value)}
                          onValueChange={(value) => field.onChange(value)}
                          min={1}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}
              {form.watch('gift_type') === 2 && (
                <FormField
                  control={form.control}
                  name='gift_value'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('giftAmount')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          placeholder={t('enterPercentage')}
                          type='number'
                          suffix='%'
                          value={field.value}
                          onValueChange={(value) => field.onChange(value)}
                          min={1}
                          max={100}
                        />
                      </FormControl>
                      <FormDescription>{t('percentageAmountDescription')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex flex-row items-center justify-end gap-2 pt-3'>
          <Button variant='outline' onClick={() => setOpen(false)}>
            {t('cancel')}
          </Button>
          <Button
            type='submit'
            form='quota-broadcast-form'
            disabled={
              isSubmitting || !form.formState.isValid || form.watch('subscribers').length === 0
            }
          >
            {isSubmitting && <Icon icon='mdi:loading' className='mr-2 h-4 w-4 animate-spin' />}
            {t('createQuotaTask')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
