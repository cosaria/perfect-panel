'use client';

import { getCurrencyConfig, updateCurrencyConfig } from '@/services/admin/system';
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

// Constants
const EXCHANGE_RATE_HOST_URL = 'https://exchangerate.host';

const currencySchema = z.object({
  access_key: z.string().optional(),
  currency_unit: z.string().min(1),
  currency_symbol: z.string().min(1),
});

type CurrencyFormData = z.infer<typeof currencySchema>;

export default function CurrencyConfig() {
  const t = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getCurrencyConfig'],
    queryFn: async () => {
      const { data } = await getCurrencyConfig();
      return data.data;
    },
    enabled: open, // Only request data when the modal is open
  });

  const form = useForm<CurrencyFormData>({
    resolver: zodResolver(currencySchema),
    defaultValues: {
      access_key: '',
      currency_unit: 'USD',
      currency_symbol: '$',
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: CurrencyFormData) {
    setLoading(true);
    try {
      await updateCurrencyConfig(values as API.CurrencyConfig);
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
              <Icon icon='mdi:currency-usd' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('currency.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('currency.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('currency.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='currency-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='access_key'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('currency.accessKey')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('currency.accessKeyPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      {t('currency.accessKeyDescription', { url: EXCHANGE_RATE_HOST_URL })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='currency_unit'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('currency.currencyUnit')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('currency.currencyUnitPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('currency.currencyUnitDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='currency_symbol'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('currency.currencySymbol')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('currency.currencySymbolPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('currency.currencySymbolDescription')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='currency-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
