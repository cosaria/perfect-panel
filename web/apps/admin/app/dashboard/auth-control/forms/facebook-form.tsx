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
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const facebookSchema = z.object({
  enabled: z.boolean(),
  client_id: z.string().optional(),
  client_secret: z.string().optional(),
});

type FacebookFormData = z.infer<typeof facebookSchema>;

export default function FacebookForm() {
  const t = useTranslations('auth-control');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch, isFetching } = useQuery({
    queryKey: ['getAuthMethodConfig', 'facebook'],
    queryFn: async () => {
      const { data } = await getAuthMethodConfig({
        method: 'facebook',
      });
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<FacebookFormData>({
    resolver: zodResolver(facebookSchema),
    defaultValues: {
      enabled: false,
      client_id: '',
      client_secret: '',
    },
  });

  useEffect(() => {
    if (data) {
      form.reset({
        enabled: data.enabled || false,
        client_id: data.config?.client_id || '',
        client_secret: data.config?.client_secret || '',
      });
    }
  }, [data, form]);

  async function onSubmit(values: FacebookFormData) {
    setLoading(true);
    try {
      await updateAuthMethodConfig({
        ...data,
        enabled: values.enabled,
        config: {
          ...data?.config,
          client_id: values.client_id,
          client_secret: values.client_secret,
        },
      } as API.UpdateAuthMethodConfigRequest);
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
              <Icon icon='mdi:facebook' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('facebook.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('facebook.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[500px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('facebook.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='facebook-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='enabled'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('facebook.enable')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('facebook.enableDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='client_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('facebook.clientId')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='1234567890123456'
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('facebook.clientIdDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='client_secret'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('facebook.clientSecret')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder='1234567890abcdef1234567890abcdef'
                        value={field.value}
                        onValueChange={field.onChange}
                        type='password'
                      />
                    </FormControl>
                    <FormDescription>{t('facebook.clientSecretDescription')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='facebook-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
