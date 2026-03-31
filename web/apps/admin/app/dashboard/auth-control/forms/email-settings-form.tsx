'use client';

import {
  getAuthMethodConfig,
  testEmailSend,
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
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@workspace/ui/components/sheet';
import { Switch } from '@workspace/ui/components/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import { Textarea } from '@workspace/ui/components/textarea';
import { HTMLEditor } from '@workspace/ui/custom-components/editor';
import { EnhancedInput } from '@workspace/ui/custom-components/enhanced-input';
import { Icon } from '@workspace/ui/custom-components/icon';
import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const emailSettingsSchema = z.object({
  id: z.number(),
  method: z.string(),
  enabled: z.boolean(),
  config: z
    .object({
      enable_verify: z.boolean(),
      enable_domain_suffix: z.boolean(),
      domain_suffix_list: z.string().optional(),
      verify_email_template: z.string().optional(),
      expiration_email_template: z.string().optional(),
      maintenance_email_template: z.string().optional(),
      traffic_exceed_email_template: z.string().optional(),
      platform: z.string(),
      platform_config: z
        .object({
          host: z.string().optional(),
          port: z.number().optional(),
          ssl: z.boolean(),
          user: z.string().optional(),
          pass: z.string().optional(),
          from: z.string().optional(),
        })
        .optional(),
    })
    .optional(),
});

type EmailSettingsFormData = z.infer<typeof emailSettingsSchema>;

export default function EmailSettingsForm() {
  const t = useTranslations('auth-control');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [testEmail, setTestEmail] = useState<string>();

  const { data, refetch, isFetching } = useQuery({
    queryKey: ['getAuthMethodConfig', 'email'],
    queryFn: async () => {
      const { data } = await getAuthMethodConfig({
        method: 'email',
      });
      return data.data;
    },
    enabled: open,
  });

  const form = useForm<EmailSettingsFormData>({
    resolver: zodResolver(emailSettingsSchema),
    defaultValues: {
      id: 0,
      method: 'email',
      enabled: false,
      config: {
        enable_verify: false,
        enable_domain_suffix: false,
        domain_suffix_list: '',
        verify_email_template: '',
        expiration_email_template: '',
        maintenance_email_template: '',
        traffic_exceed_email_template: '',
        platform: 'smtp',
        platform_config: {
          host: '',
          port: 587,
          ssl: false,
          user: '',
          pass: '',
          from: '',
        },
      },
    },
  });

  useEffect(() => {
    if (data) {
      form.reset(data);
    }
  }, [data, form]);

  async function onSubmit(values: EmailSettingsFormData) {
    setLoading(true);
    try {
      await updateAuthMethodConfig({
        ...values,
        config: {
          ...values.config,
          platform: 'smtp',
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
              <Icon icon='mdi:email-outline' className='text-primary h-5 w-5' />
            </div>
            <div className='flex-1'>
              <p className='font-medium'>{t('email.title')}</p>
              <p className='text-muted-foreground text-sm'>{t('email.description')}</p>
            </div>
          </div>
          <Icon icon='mdi:chevron-right' className='size-6' />
        </div>
      </SheetTrigger>
      <SheetContent className='w-[600px] max-w-full md:max-w-screen-md'>
        <SheetHeader>
          <SheetTitle>{t('email.title')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className='-mx-6 h-[calc(100dvh-48px-36px-36px-env(safe-area-inset-top))] px-6'>
          <Form {...form}>
            <form
              id='email-settings-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-2 pt-4'
            >
              <Tabs defaultValue='basic' className='space-y-2'>
                <TabsList className='flex h-full w-full flex-wrap *:flex-auto md:flex-nowrap'>
                  <TabsTrigger value='basic'>{t('email.basicSettings')}</TabsTrigger>
                  <TabsTrigger value='smtp'>{t('email.smtpSettings')}</TabsTrigger>
                  <TabsTrigger value='verify'>{t('email.verifyTemplate')}</TabsTrigger>
                  <TabsTrigger value='expiration'>{t('email.expirationTemplate')}</TabsTrigger>
                  <TabsTrigger value='maintenance'>{t('email.maintenanceTemplate')}</TabsTrigger>
                  <TabsTrigger value='traffic'>{t('email.trafficTemplate')}</TabsTrigger>
                </TabsList>

                <TabsContent value='basic' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='enabled'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.enable')}</FormLabel>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            className='float-end !mt-0'
                          />
                        </FormControl>
                        <FormDescription>{t('email.enableDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.enable_verify'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.emailVerification')}</FormLabel>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            className='float-end !mt-0'
                          />
                        </FormControl>
                        <FormDescription>{t('email.emailVerificationDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.enable_domain_suffix'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.emailSuffixWhitelist')}</FormLabel>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            className='float-end !mt-0'
                          />
                        </FormControl>
                        <FormDescription>
                          {t('email.emailSuffixWhitelistDescription')}
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.domain_suffix_list'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.whitelistSuffixes')}</FormLabel>
                        <FormControl>
                          <Textarea
                            className='h-32'
                            placeholder={t('email.whitelistSuffixesPlaceholder')}
                            value={field.value}
                            onChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>{t('email.whitelistSuffixesDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>

                <TabsContent value='smtp' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='config.platform_config.host'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.smtpServerAddress')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onValueChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>{t('email.smtpServerAddressDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.platform_config.port'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.smtpServerPort')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            type='number'
                            placeholder='587'
                            value={field.value?.toString()}
                            onValueChange={(value) => field.onChange(Number(value))}
                          />
                        </FormControl>
                        <FormDescription>{t('email.smtpServerPortDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.platform_config.ssl'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.smtpEncryptionMethod')}</FormLabel>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            className='float-end !mt-0'
                          />
                        </FormControl>
                        <FormDescription>
                          {t('email.smtpEncryptionMethodDescription')}
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.platform_config.user'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.smtpAccount')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onValueChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>{t('email.smtpAccountDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.platform_config.pass'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.smtpPassword')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            type='password'
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onValueChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>{t('email.smtpPasswordDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name='config.platform_config.from'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.senderAddress')}</FormLabel>
                        <FormControl>
                          <EnhancedInput
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onValueChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>{t('email.senderAddressDescription')}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className='space-y-2 border-t pt-4'>
                    <FormLabel>{t('email.sendTestEmail')}</FormLabel>
                    <div className='flex items-center gap-2'>
                      <EnhancedInput
                        placeholder='test@example.com'
                        type='email'
                        value={testEmail}
                        onValueChange={(value) => setTestEmail(value as string)}
                      />
                      <Button
                        type='button'
                        disabled={!testEmail || isFetching}
                        onClick={async () => {
                          if (!testEmail) return;
                          try {
                            await testEmailSend({ email: testEmail });
                            toast.success(t('email.sendSuccess'));
                          } catch {
                            toast.error(t('email.sendFailure'));
                          }
                        }}
                      >
                        {t('email.sendTestEmail')}
                      </Button>
                    </div>
                    <p className='text-muted-foreground text-xs'>
                      {t('email.sendTestEmailDescription')}
                    </p>
                  </div>
                </TabsContent>

                <TabsContent value='verify' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='config.verify_email_template'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.verifyEmailTemplate')}</FormLabel>
                        <FormControl>
                          <HTMLEditor
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onBlur={field.onChange}
                          />
                        </FormControl>
                        <div className='mt-4 space-y-2 border-t pt-4'>
                          <p className='text-muted-foreground text-sm font-medium'>
                            {t('email.templateVariables.title')}
                          </p>
                          <div className='text-muted-foreground space-y-2 text-xs'>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.Type}}'}
                              </code>
                              <span>{t('email.templateVariables.type.description')}</span>
                            </div>
                            <div className='pl-6 text-orange-600 dark:text-orange-400'>
                              💡 {t('email.templateVariables.type.conditionalSyntax')}
                              <br />
                              <code className='rounded bg-orange-50 px-1 text-xs dark:bg-orange-900/20'>
                                {'{{if eq .Type 1}}...{{else}}...{{end}}'}
                              </code>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteLogo}}'}
                              </code>
                              <span>{t('email.templateVariables.siteLogo.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteName}}'}
                              </code>
                              <span>{t('email.templateVariables.siteName.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.Expire}}'}
                              </code>
                              <span>{t('email.templateVariables.expire.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.Code}}'}
                              </code>
                              <span>{t('email.templateVariables.code.description')}</span>
                            </div>
                          </div>
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>

                <TabsContent value='expiration' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='config.expiration_email_template'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.expirationEmailTemplate')}</FormLabel>
                        <FormControl>
                          <HTMLEditor
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onBlur={field.onChange}
                          />
                        </FormControl>
                        <div className='mt-4 space-y-2 border-t pt-4'>
                          <p className='text-muted-foreground text-sm font-medium'>
                            {t('email.templateVariables.title')}
                          </p>
                          <div className='text-muted-foreground space-y-2 text-xs'>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteLogo}}'}
                              </code>
                              <span>{t('email.templateVariables.siteLogo.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteName}}'}
                              </code>
                              <span>{t('email.templateVariables.siteName.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.ExpireDate}}'}
                              </code>
                              <span>{t('email.templateVariables.expireDate.description')}</span>
                            </div>
                          </div>
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>

                <TabsContent value='maintenance' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='config.maintenance_email_template'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.maintenanceEmailTemplate')}</FormLabel>
                        <FormControl>
                          <HTMLEditor
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onBlur={field.onChange}
                          />
                        </FormControl>
                        <div className='mt-4 space-y-2 border-t pt-4'>
                          <p className='text-muted-foreground text-sm font-medium'>
                            {t('email.templateVariables.title')}
                          </p>
                          <div className='text-muted-foreground space-y-2 text-xs'>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteLogo}}'}
                              </code>
                              <span>{t('email.templateVariables.siteLogo.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteName}}'}
                              </code>
                              <span>{t('email.templateVariables.siteName.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.MaintenanceDate}}'}
                              </code>
                              <span>
                                {t('email.templateVariables.maintenanceDate.description')}
                              </span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.MaintenanceTime}}'}
                              </code>
                              <span>
                                {t('email.templateVariables.maintenanceTime.description')}
                              </span>
                            </div>
                          </div>
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>

                <TabsContent value='traffic' className='space-y-2'>
                  <FormField
                    control={form.control}
                    name='config.traffic_exceed_email_template'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email.trafficExceedEmailTemplate')}</FormLabel>
                        <FormControl>
                          <HTMLEditor
                            placeholder={t('email.inputPlaceholder')}
                            value={field.value}
                            onBlur={field.onChange}
                          />
                        </FormControl>
                        <div className='mt-4 space-y-2 border-t pt-4'>
                          <p className='text-muted-foreground text-sm font-medium'>
                            {t('email.templateVariables.title')}
                          </p>
                          <div className='text-muted-foreground space-y-2 text-xs'>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteLogo}}'}
                              </code>
                              <span>{t('email.templateVariables.siteLogo.description')}</span>
                            </div>
                            <div className='flex items-center gap-2'>
                              <code className='bg-muted text-foreground rounded px-1.5 py-0.5 font-mono'>
                                {'{{.SiteName}}'}
                              </code>
                              <span>{t('email.templateVariables.siteName.description')}</span>
                            </div>
                          </div>
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>
              </Tabs>
            </form>
          </Form>
        </ScrollArea>
        <SheetFooter className='flex-row justify-end gap-2 pt-3'>
          <Button variant='outline' disabled={loading} onClick={() => setOpen(false)}>
            {t('common.cancel')}
          </Button>
          <Button disabled={loading} type='submit' form='email-settings-form'>
            {loading && <Icon icon='mdi:loading' className='mr-2 animate-spin' />}
            {t('common.save')}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
