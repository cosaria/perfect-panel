'use client';

import { getTosConfig, updateTosConfig } from '@/services/admin/system';
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
import { MarkdownEditor } from '@workspace/ui/custom-components/editor';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const tosSchema = z.object({
  tos_content: z.string().optional(),
});

type TosFormData = z.infer<typeof tosSchema>;

export default function TosConfig() {
  const t = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getTosConfig'],
    queryFn: async () => {
      const { data } = await getTosConfig();
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<TosFormData>({
    resolver: zodResolver(tosSchema),
    defaultValues: {
      tos_content: '',
    },
  });

  useEffect(() => {
    if (data) {
      form.reset({
        tos_content: data.tos_content || '',
      });
    }
  }, [data, form]);

  async function onSubmit(values: TosFormData) {
    setLoading(true);
    try {
      await updateTosConfig(values as API.TosConfig);
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
              <Icon icon='mdi:file-document-outline' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('tos.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('tos.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('tos.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form id='tos-form' onSubmit={form.handleSubmit(onSubmit)} className='space-y-2 pt-4'>
              <FormField
                control={form.control}
                name='tos_content'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('tos.title')}</FormLabel>
                    <FormControl>
                      <MarkdownEditor value={field.value || ''} onChange={field.onChange} />
                    </FormControl>
                    <FormDescription>{t('tos.description')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='tos-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
