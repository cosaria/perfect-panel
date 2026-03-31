'use client';

import { getSiteConfig, updateSiteConfig } from '@/services/admin/system';
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

import { Textarea } from '@workspace/ui/components/textarea';
import { JSONEditor } from '@workspace/ui/custom-components/editor';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { UploadImage } from '@workspace/ui/custom-components/upload-image';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const siteSchema = z.object({
  site_logo: z.string().optional(),
  site_name: z.string().min(1),
  site_desc: z.string().optional(),
  keywords: z.string().optional(),
  custom_html: z.string().optional(),
  host: z.string().optional(),
  custom_data: z.any().optional(),
});

type SiteFormData = z.infer<typeof siteSchema>;

export default function SiteConfig() {
  const t = useTranslations('system');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const { data, refetch } = useQuery({
    queryKey: ['getSiteConfig'],
    queryFn: async () => {
      const { data } = await getSiteConfig();
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<SiteFormData>({
    resolver: zodResolver(siteSchema),
    defaultValues: {
      site_logo: '',
      site_name: '',
      site_desc: '',
      keywords: '',
      custom_html: '',
      host: '',
      custom_data: {},
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: SiteFormData) {
    setLoading(true);
    try {
      await updateSiteConfig(values as API.SiteConfig);
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
              <Icon icon='mdi:web' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('site.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('site.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('site.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form id='site-form' onSubmit={form.handleSubmit(onSubmit)} className='space-y-2 pt-4'>
              <FormField
                control={form.control}
                name='site_logo'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.logo')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('site.logoPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                        suffix={
                          <UploadImage
                            className='bg-muted h-9 rounded-none border-none px-2'
                            onChange={(value) => {
                              field.onChange(value);
                            }}
                          />
                        }
                      />
                    </FormControl>
                    <FormDescription>{t('site.logoDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='site_name'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.siteName')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('site.siteNamePlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('site.siteNameDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='site_desc'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.siteDesc')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('site.siteDescPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('site.siteDescDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='keywords'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.keywords')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        placeholder={t('site.keywordsPlaceholder')}
                        value={field.value}
                        onValueChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('site.keywordsDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='custom_html'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.customHtml')}</FormLabel>
                    <FormControl>
                      <Textarea
                        className='h-32'
                        placeholder={t('site.customHtmlDescription')}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>{t('site.customHtmlDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='host'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.siteDomain')}</FormLabel>
                    <FormControl>
                      <Textarea
                        className='h-32'
                        placeholder={`${t('site.siteDomainPlaceholder')}\nexample.com\nwww.example.com`}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>{t('site.siteDomainDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='custom_data'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('site.customData')}</FormLabel>
                    <FormControl>
                      <JSONEditor
                        schema={{
                          type: 'object',
                          additionalProperties: true,
                          properties: {
                            website: { type: 'string', title: 'Website' },
                            contacts: {
                              type: 'object',
                              title: 'Contacts',
                              additionalProperties: true,
                              properties: {
                                email: { type: 'string', title: 'Email' },
                                telephone: { type: 'string', title: 'Telephone' },
                                address: { type: 'string', title: 'Address' },
                              },
                            },
                            community: {
                              type: 'object',
                              title: 'Community',
                              additionalProperties: true,
                              properties: {
                                telegram: { type: 'string', title: 'Telegram' },
                                twitter: { type: 'string', title: 'Twitter' },
                                discord: { type: 'string', title: 'Discord' },
                                instagram: { type: 'string', title: 'Instagram' },
                                linkedin: { type: 'string', title: 'Linkedin' },
                                facebook: { type: 'string', title: 'Facebook' },
                                github: { type: 'string', title: 'Github' },
                              },
                            },
                          },
                        }}
                        value={field.value}
                        onBlur={(value) => field.onChange(value)}
                      />
                    </FormControl>
                    <FormDescription>{t('site.customDataDescription')}</FormDescription>
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
          <Button disabled={loading} type='submit' form='site-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
