'use client';

import {
  getAuthMethodConfig,
  getSmsPlatform,
  testSmsSend,
  updateAuthMethodConfig,
} from '@/services/admin/authMethod';
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
import { Switch } from '@workspace/ui/components/switch';
import { Textarea } from '@workspace/ui/components/textarea';
import { AreaCodeSelect } from '@workspace/ui/custom-components/area-code-select';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import TagInput from '@workspace/ui/custom-components/tag-input';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const phoneSettingsSchema = z.object({
  id: z.number(),
  method: z.string(),
  enabled: z.boolean(),
  config: z
    .object({
      enable_whitelist: z.boolean().optional(),
      whitelist: z.array(z.string()).optional(),
      platform: z.string().optional(),
      platform_config: z
        .object({
          access: z.string().optional(),
          endpoint: z.string().optional(),
          secret: z.string().optional(),
          template_code: z.string().optional(),
          sign_name: z.string().optional(),
          phone_number: z.string().optional(),
          template: z.string().optional(),
        })
        .optional(),
    })
    .optional(),
});

type PhoneSettingsFormData = z.infer<typeof phoneSettingsSchema>;

export default function PhoneSettingsForm() {
  const t = useTranslations('auth-control');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [testParams, setTestParams] = useState<API.TestSmsSendRequest>({
    telephone: '',
    area_code: '1',
  });

  const { data, refetch, isFetching } = useQuery({
    queryKey: ['getAuthMethodConfig', 'mobile'],
    queryFn: async () => {
      const { data } = await getAuthMethodConfig({
        method: 'mobile',
      });
      return data.data;
    },
    enabled: open,
  });

  const { data: platforms } = useQuery({
    queryKey: ['getSmsPlatform'],
    queryFn: async () => {
      const { data } = await getSmsPlatform();
      return data.data?.list;
    },
    enabled: open,
  });

  const form = useForm<PhoneSettingsFormData>({
    resolver: zodResolver(phoneSettingsSchema),
    defaultValues: {
      id: 0,
      method: 'mobile',
      enabled: false,
      config: {
        enable_whitelist: false,
        whitelist: [],
        platform: '',
        platform_config: {
          access: '',
          endpoint: '',
          secret: '',
          template_code: 'code',
          sign_name: '',
          phone_number: '',
          template: '',
        },
      },
    },
  });

  const selectedPlatform = platforms?.find(
    (platform) => platform.platform === form.watch('config.platform'),
  );
  const { platform_url, platform_field_description: platformConfig } = selectedPlatform ?? {};

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: PhoneSettingsFormData) {
    setLoading(true);
    try {
      await updateAuthMethodConfig(values as API.UpdateAuthMethodConfigRequest);
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
              <Icon icon='mdi:phone-settings' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('phone.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('phone.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('phone.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='phone-settings-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <FormField
                control={form.control}
                name='enabled'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('phone.enable')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        disabled={isFetching}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('phone.enableTip')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.enable_whitelist'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('phone.whitelistValidation')}</FormLabel>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        className='float-end !mt-0'
                      />
                    </FormControl>
                    <FormDescription>{t('phone.whitelistValidationTip')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.whitelist'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('phone.whitelistAreaCode')}</FormLabel>
                    <FormControl>
                      <TagInput
                        placeholder='1, 852, 886, 888'
                        value={field.value}
                        onChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>{t('phone.whitelistAreaCodeTip')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.platform'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('phone.platform')}</FormLabel>
                    <div className='flex items-center gap-1'>
                      <FormControl>
                        <Select
                          value={field.value}
                          onValueChange={field.onChange}
                          disabled={isFetching}
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {platforms?.map((item) => (
                              <SelectItem key={item.platform} value={item.platform}>
                                {item.platform}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </FormControl>
                      {platform_url && (
                        <Button size='sm' asChild>
                          <Link href={platform_url} target='_blank'>
                            {t('phone.applyPlatform')}
                          </Link>
                        </Button>
                      )}
                    </div>
                    <FormDescription>{t('phone.platformTip')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='config.platform_config.access'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('phone.accessLabel')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        value={field.value}
                        onValueChange={field.onChange}
                        disabled={isFetching}
                        placeholder={t('phone.platformConfigTip', { key: platformConfig?.access })}
                      />
                    </FormControl>
                    <FormDescription>
                      {t('phone.platformConfigTip', { key: platformConfig?.access })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {platformConfig?.endpoint && (
                <FormField
                  control={form.control}
                  name='config.platform_config.endpoint'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('phone.endpointLabel')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          value={field.value}
                          onValueChange={field.onChange}
                          disabled={isFetching}
                          placeholder={t('phone.platformConfigTip', {
                            key: platformConfig?.endpoint,
                          })}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('phone.platformConfigTip', { key: platformConfig?.endpoint })}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              <FormField
                control={form.control}
                name='config.platform_config.secret'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('phone.secretLabel')}</FormLabel>
                    <FormControl>
                      <EnhancedInput
                        type='password'
                        value={field.value}
                        onValueChange={field.onChange}
                        disabled={isFetching}
                        placeholder={t('phone.platformConfigTip', { key: platformConfig?.secret })}
                      />
                    </FormControl>
                    <FormDescription>
                      {t('phone.platformConfigTip', { key: platformConfig?.secret })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {platformConfig?.template_code && (
                <FormField
                  control={form.control}
                  name='config.platform_config.template_code'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('phone.templateCodeLabel')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          value={field.value}
                          onValueChange={field.onChange}
                          disabled={isFetching}
                          placeholder={t('phone.platformConfigTip', {
                            key: platformConfig?.template_code,
                          })}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('phone.platformConfigTip', { key: platformConfig?.template_code })}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              {platformConfig?.sign_name && (
                <FormField
                  control={form.control}
                  name='config.platform_config.sign_name'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('phone.signNameLabel')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          value={field.value}
                          onValueChange={field.onChange}
                          disabled={isFetching}
                          placeholder={t('phone.platformConfigTip', {
                            key: platformConfig?.sign_name,
                          })}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('phone.platformConfigTip', { key: platformConfig?.sign_name })}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              {platformConfig?.phone_number && (
                <FormField
                  control={form.control}
                  name='config.platform_config.phone_number'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('phone.phoneNumberLabel')}</FormLabel>
                      <FormControl>
                        <EnhancedInput
                          value={field.value}
                          onValueChange={field.onChange}
                          disabled={isFetching}
                          placeholder={t('phone.platformConfigTip', {
                            key: platformConfig?.phone_number,
                          })}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('phone.platformConfigTip', { key: platformConfig?.phone_number })}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              {platformConfig?.code_variable && (
                <FormField
                  control={form.control}
                  name='config.platform_config.template'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('phone.template')}</FormLabel>
                      <FormControl>
                        <Textarea
                          value={field.value}
                          onChange={field.onChange}
                          disabled={isFetching}
                          placeholder={t('phone.placeholders.template', {
                            code: platformConfig?.code_variable,
                          })}
                        />
                      </FormControl>
                      <FormDescription>
                        {t('phone.templateTip', { code: platformConfig?.code_variable })}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              <div className='space-y-4 border-t pt-4'>
                <div>
                  <FormLabel>{t('phone.testSms')}</FormLabel>
                  <p className='text-muted-foreground mb-3 text-sm'>{t('phone.testSmsTip')}</p>
                  <div className='flex items-center gap-2'>
                    <AreaCodeSelect
                      value={testParams.area_code}
                      onChange={(value) => {
                        if (value.phone) {
                          setTestParams((prev) => ({ ...prev, area_code: value.phone! }));
                        }
                      }}
                    />
                    <EnhancedInput
                      placeholder={t('phone.testSmsPhone')}
                      value={testParams.telephone}
                      onValueChange={(value) => {
                        setTestParams((prev) => ({ ...prev, telephone: value as string }));
                      }}
                    />
                    <Button
                      type='button'
                      disabled={!testParams.telephone || !testParams.area_code || isFetching}
                      onClick={async () => {
                        if (isFetching || !testParams.telephone || !testParams.area_code) return;
                        try {
                          await testSmsSend(testParams);
                          toast.success(t('phone.sendSuccess'));
                        } catch {
                          toast.error(t('phone.sendFailed'));
                        }
                      }}
                    >
                      {t('phone.testSms')}
                    </Button>
                  </div>
                </div>
              </div>
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {t('common.cancel')}
          </Button>
          <Button disabled={loading} type='submit' form='phone-settings-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
