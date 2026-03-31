'use client';

import { getInviteConfig, updateInviteConfig } from '@/services/admin/system';
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

const inviteSchema = z.object({
  forced_invite: z.boolean().optional(),
  referral_percentage: z.number().optional(),
  only_first_purchase: z.boolean().optional(),
});

type InviteFormData = z.infer<typeof inviteSchema>;

export default function InviteConfig() {
  const t = useTranslations('system.invite');
  const systemT = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getInviteConfig'],
    queryFn: async () => {
      const { data } = await getInviteConfig();
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<InviteFormData>({
    resolver: zodResolver(inviteSchema),
    defaultValues: {
      forced_invite: false,
      referral_percentage: 0,
      only_first_purchase: false,
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: InviteFormData) {
    setLoading(true);
    try {
      await updateInviteConfig(values as API.InviteConfig);
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
              <Icon icon='mdi:account-multiple-plus-outline' className='text-primary h-5 w-5' />
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
              id='invite-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='forced_invite'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('forcedInvite')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('forcedInviteDescription')}</FormDescription>
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
                        placeholder={t('inputPlaceholder')}
                        value={field.value}
                        type='number'
                        min={0}
                        max={100}
                        suffix='%'
                        onValueBlur={(value) => field.onChange(Number(value))}
                      />
                    </FormControl>
                    <FormDescription>{t('referralPercentageDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='only_first_purchase'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('onlyFirstPurchase')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('onlyFirstPurchaseDescription')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='invite-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {systemT('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
